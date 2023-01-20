package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/repoitory"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GopherMartOrderController struct {
	repository repoitory.OrderTable
}

var ErrInvalidOrderId = errors.New("invalid order id")
var ErrExistingOrderForCurrentUser = errors.New("order existing for current user")
var ErrExistingOrderForAnotherUser = errors.New("order existing for another user")

func InitOrderController(table repoitory.OrderTable) *GopherMartOrderController {
	return &GopherMartOrderController{repository: table}
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
	return c.repository.GetAllOrdersByLogin(login, context)
}

func (c *GopherMartOrderController) checkOrder(login, orderId string, userController *GopherMartUserController) {
	for {
		time.Sleep(time.Millisecond * 500)
		r, err := http.Get("lol")
		if err != nil {
			c.repository.UpdateOrder(orderId, types.InvalidOrder, 0, context.Background())
			return
		}
		d := json.NewDecoder(r.Body)
		resp := new(types.LoyaltyServiceResponse)
		err = d.Decode(resp)
		if err != nil {
			c.repository.UpdateOrder(orderId, types.InvalidOrder, 0, context.Background())
			fmt.Println(err)
		}
		if r.Body.Close() != nil {
			fmt.Println(err)
		}
		switch resp.Status {
		case types.LoyaltyServiceRegistered:
			continue
		case types.LoyaltyServiceProcessing:
			c.repository.UpdateOrder(orderId, types.ProcessingOrder, 0, context.Background())
			continue
		case types.LoyaltyServiceProcessed:
			c.repository.UpdateOrder(orderId, types.ProcessedOrder, resp.Accrual, context.Background())
			userController.UpdateUserBalance(login, resp.Accrual, context.Background())
			return
		case types.LoyaltyServiceInvalid:
			c.repository.UpdateOrder(orderId, types.InvalidOrder, 0, context.Background())
			return
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
