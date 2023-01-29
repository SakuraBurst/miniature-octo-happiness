package repoitory

import (
	"context"
	"fmt"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"time"
)

type UserTable interface {
	CreateUser(login, hashedPassword string, c context.Context) error
	GetUser(login string, c context.Context) (*types.User, error)
	UpdateBalanceAndWithdraw(login string, newBalance, newWithdraw float64, c context.Context) error
}

type OrderTable interface {
	CreateOrder(login, orderID string, c context.Context) error
	UpdateOrder(orderID string, status types.OrderStatus, accrual float64, c context.Context) error
	GetOrderByOrderID(orderID string, c context.Context) (*types.Order, error)
	GetAllOrdersByLogin(login string, c context.Context) ([]types.Order, error)
}

type WithdrawTable interface {
	CreateWithdraw(login, orderID string, sum float64, c context.Context) error
	GetAllWithdrawalsByLogin(login string, c context.Context) ([]types.Withdraw, error)
}

func InitDataBase(address string) (UserTable, OrderTable, WithdrawTable) {
	configSql, err := os.ReadFile("cmd/gophermart/config/init.sql")
	if err != nil {
		log.Fatal(err)
	}
	c, cf := context.WithTimeout(context.Background(), time.Second)
	defer cf()
	conn, err := pgx.Connect(c, address)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	_, err = conn.Exec(c, string(configSql))
	if err != nil {
		log.Fatal(err)
	}
	return &userTable{conn}, &ordersTable{conn}, &withdrawTable{conn}
}
