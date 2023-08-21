package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/repository"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GopherMartOrderController struct {
	mu                        *sync.RWMutex
	repository                repository.OrderTable
	loyaltyServiceBaseAddress *url.URL
	checkOrderQueue           chan struct{}
	isQueueClosed             bool
}

var ErrInvalidOrderID = errors.New("invalid order id")
var ErrExistingOrderForCurrentUser = errors.New("order existing for current user")
var ErrExistingOrderForAnotherUser = errors.New("order existing for another user")
var ErrShutdown = errors.New("server is shutting down")

func InitOrderController(table repository.OrderTable, loyaltyServiceBaseAddress string) (*GopherMartOrderController, error) {
	u, err := url.Parse(loyaltyServiceBaseAddress)
	u.Path = "/api/orders/"
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &GopherMartOrderController{
		mu:                        &sync.RWMutex{},
		repository:                table,
		loyaltyServiceBaseAddress: u,
		checkOrderQueue:           make(chan struct{}, 10),
		isQueueClosed:             false,
	}, nil
}

func (c *GopherMartOrderController) CreateOrder(orderID, login string, userController *GopherMartUserController, context context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.isQueueClosed {
		return ErrShutdown
	}
	if !Luhn(orderID) {
		return ErrInvalidOrderID
	}
	order, err := c.repository.GetOrderByOrderID(orderID, context)
	if err == nil {
		if order.UserLogin == login {
			return ErrExistingOrderForCurrentUser
		} else {
			return ErrExistingOrderForAnotherUser
		}
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		log.Error(err)
		return err
	}
	err = c.repository.CreateOrder(login, orderID, context)
	if err != nil {

		return err
	}
	c.checkOrderQueue <- struct{}{}
	go c.checkOrder(login, orderID, userController)
	return nil
}

func (c *GopherMartOrderController) CloseQueue() {
	c.mu.Lock()
	defer c.mu.Unlock()
	close(c.checkOrderQueue)
	c.isQueueClosed = true
}

func (c *GopherMartOrderController) IsQueueEmpty() bool {
	return len(c.checkOrderQueue) == 0
}

func (c *GopherMartOrderController) GetUserOrders(login string, context context.Context) ([]types.Order, error) {
	o, err := c.repository.GetAllOrdersByLogin(login, context)
	if err == nil {
		log.Error(err)
		return o, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return nil, err
}

func (c *GopherMartOrderController) checkOrder(login, orderID string, userController *GopherMartUserController) {
	orderURL := c.loyaltyServiceBaseAddress.String() + orderID
	defer func() {
		<-c.checkOrderQueue
	}()
	for {
		time.Sleep(time.Millisecond * 100)
		r, err := http.Get(orderURL)
		if err != nil {
			err = c.repository.UpdateOrder(orderID, types.InvalidOrder, 0, context.Background())
			if err != nil {
				log.Error(err)
			}
			return
		}
		resp := new(types.LoyaltyServiceResponse)
		err = json.NewDecoder(r.Body).Decode(resp)
		if err != nil {
			log.Error(err)
			err = c.repository.UpdateOrder(orderID, types.InvalidOrder, 0, context.Background())
			if err != nil {
				log.Println(err)
			}
			return
		}
		if r.Body.Close() != nil {
			log.Error(err)
			return
		}
		switch resp.Status {
		case types.LoyaltyServiceRegistered:
			continue
		case types.LoyaltyServiceProcessing:
			err = c.repository.UpdateOrder(orderID, types.ProcessingOrder, 0, context.Background())
			if err != nil {
				log.Error(err)
				break
			}
			continue
		case types.LoyaltyServiceProcessed:
			err = c.repository.UpdateOrder(orderID, types.ProcessedOrder, resp.Accrual, context.Background())
			if err != nil {
				log.Error(err)
				break
			}
			err = userController.AddUserBalance(login, resp.Accrual, context.Background())
			if err != nil {
				log.Error(err)
			}
			return
		case types.LoyaltyServiceInvalid:
			err = c.repository.UpdateOrder(orderID, types.InvalidOrder, 0, context.Background())
			if err != nil {
				log.Error(err)
			}
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
