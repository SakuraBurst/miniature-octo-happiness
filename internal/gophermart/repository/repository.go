package repository

import (
	"context"
	"database/sql/driver"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
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

type DB interface {
	driver.Pinger
	Close(ctx context.Context) error
}

func InitDataBase(address string) (UserTable, OrderTable, WithdrawTable, DB, error) {
	configSql, err := os.ReadFile("cmd/gophermart/config/init.sql")
	if err != nil {
		return nil, nil, nil, nil, err
	}
	c, cf := context.WithTimeout(context.Background(), time.Second)
	defer cf()
	conn, err := pgx.Connect(c, address)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	_, err = conn.Exec(c, string(configSql))
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return &userTable{conn}, &ordersTable{conn}, &withdrawTable{conn}, conn, nil
}
