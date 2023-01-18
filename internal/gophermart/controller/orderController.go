package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/repoitory"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
	"net/http"
	"strconv"
	"strings"
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

func (c *GopherMartOrderController) CreateOrder(orderId, login string, context context.Context) error {
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
		fmt.Println(err)
		return err
	}
	err = c.repository.CreateOrder(login, orderId, context)
	if err != nil {
		return err
	}
	go c.checkOrder(orderId)
	return nil
}

func (c *GopherMartOrderController) GetUserOrders(login string, context context.Context) ([]types.Order, error) {
	return c.repository.GetAllOrdersByLogin(login, context)
}

func (c *GopherMartOrderController) checkOrder(orderId string) {
	for {
		r, err := http.Get("lol")
		if err != nil {
			c.repository.UpdateOrder(orderId, types.Invalid, 0, context.Background())
			return
		}
		fmt.Println(r)
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
