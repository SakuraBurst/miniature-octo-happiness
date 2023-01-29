package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/repoitory"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type GopherMartOrderController struct {
	repository                repoitory.OrderTable
	loyaltyServiceBaseAddress *url.URL
}

var ErrInvalidOrderId = errors.New("invalid order id")
var ErrExistingOrderForCurrentUser = errors.New("order existing for current user")
var ErrExistingOrderForAnotherUser = errors.New("order existing for another user")

func InitOrderController(table repoitory.OrderTable, loyaltyServiceBaseAddress string) *GopherMartOrderController {
	fmt.Println(loyaltyServiceBaseAddress)
	u, err := url.Parse(loyaltyServiceBaseAddress)
	if err != nil {
		log.Fatal(err)
	}
	return &GopherMartOrderController{repository: table, loyaltyServiceBaseAddress: u}
}

func (c *GopherMartOrderController) CreateOrder(orderId, login string, userController *GopherMartUserController, context context.Context) error {
	if !Luhn(orderId) {
		return ErrInvalidOrderId
	}
	order, err := c.repository.GetOrderByOrderId(orderId, context)
	if err == nil {
		if order.UserLogin == login {
			return ErrExistingOrderForCurrentUser
		} else {
			return ErrExistingOrderForAnotherUser
		}
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	err = c.repository.CreateOrder(login, orderId, context)
	if err != nil {

		return err
	}
	go c.checkOrder(login, orderId, userController)
	return nil
}

func (c *GopherMartOrderController) GetUserOrders(login string, context context.Context) ([]types.Order, error) {
	o, err := c.repository.GetAllOrdersByLogin(login, context)
	if err == nil {
		return o, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return nil, err
}

func (c *GopherMartOrderController) checkOrder(login, orderId string, userController *GopherMartUserController) {
	c.loyaltyServiceBaseAddress.Path = "/api/orders/" + orderId
	defer func() {
		c.loyaltyServiceBaseAddress.Path = ""
	}()
	for {
		time.Sleep(time.Millisecond * 500)
		r, err := http.Get(c.loyaltyServiceBaseAddress.String())
		if err != nil {
			err = c.repository.UpdateOrder(orderId, types.InvalidOrder, 0, context.Background())
			if err != nil {
				log.Println(err)
			}
			return
		}
		resp := new(types.LoyaltyServiceResponse)
		err = json.NewDecoder(r.Body).Decode(resp)
		if err != nil {
			log.Println(err)
			err = c.repository.UpdateOrder(orderId, types.InvalidOrder, 0, context.Background())
			if err != nil {
				log.Println(err)
			}
			return
		}
		if r.Body.Close() != nil {
			log.Println(err)
			return
		}
		switch resp.Status {
		case types.LoyaltyServiceRegistered:
			continue
		case types.LoyaltyServiceProcessing:
			err = c.repository.UpdateOrder(orderId, types.ProcessingOrder, 0, context.Background())
			if err != nil {
				break
			}
			continue
		case types.LoyaltyServiceProcessed:
			err = c.repository.UpdateOrder(orderId, types.ProcessedOrder, resp.Accrual, context.Background())
			if err != nil {
				log.Println(err)
			}
			err = userController.UpdateUserBalance(login, resp.Accrual, context.Background())
			if err != nil {
				log.Println(err)
			}
			return
		case types.LoyaltyServiceInvalid:
			err = c.repository.UpdateOrder(orderId, types.InvalidOrder, 0, context.Background())
			if err != nil {
				log.Println(err)
			}
			return
		}
		if err != nil {
			break
		}
	}
}

func Luhn(val string) bool {
	val = strings.Replace(val, " ", "", -1)
	sum, err := strconv.Atoi(val[len(val)-1:])
	if err != nil {
		return false
	}
	lastIndex := len(val) - 1
	parity := len(val) % 2
	for i, v := range val {
		if i == lastIndex {
			break
		}
		c, err := strconv.Atoi(string(v))
		if err != nil {
			return false
		}
		if i%2 == parity {
			prod := c << 1
			if prod > 9 {
				c = prod - 9
			} else {
				c = prod
			}
		}
		sum += c
	}
	return sum%10 == 0
}
