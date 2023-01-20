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

type userTable struct {
	*pgx.Conn
}

type UserTable interface {
	CreateUser(login, hashedPassword string, c context.Context) error
	GetUser(login string, c context.Context) (*types.User, error)
	UpdateBalanceAndWithdraw(login string, newBalance, newWithdraw float64, c context.Context) error
}

func (ut *userTable) CreateUser(login, hashedPassword string, c context.Context) error {
	_, err := ut.Exec(c, "insert into users(login, password) values ($1, $2)", login, hashedPassword)
	return err
}

func (ut *userTable) GetUser(login string, c context.Context) (*types.User, error) {
	r := ut.QueryRow(c, "select * from users where login = $1", login)
	user := new(types.User)
	err := r.Scan(&user.Id, &user.Login, &user.Password, &user.Balance, &user.Withdraw)
	return user, err
}

func (ut *userTable) UpdateBalanceAndWithdraw(login string, newBalance, newWithdraw float64, c context.Context) error {
	_, err := ut.Exec(c, "update users set balance = $2, withdraw = $3 where login = $1", login, newBalance, newWithdraw)
	return err
}

type ordersTable struct {
	*pgx.Conn
}

type OrderTable interface {
	CreateOrder(login, orderId string, c context.Context) error
	UpdateOrder(orderId string, status types.OrderStatus, accrual float64, c context.Context) error
	GetOrderByOrderId(orderId string, c context.Context) (*types.Order, error)
	GetAllOrdersByLogin(login string, c context.Context) ([]types.Order, error)
}

func (ot *ordersTable) CreateOrder(login, orderId string, c context.Context) error {
	_, err := ot.Exec(c, "insert into orders(status, user_login, order_id, uploaded_at) values ( 'NEW',$1, $2, $3)", login, orderId, time.Now())
	return err
}

func (ot *ordersTable) UpdateOrder(orderId string, status types.OrderStatus, accrual float64, c context.Context) error {
	_, err := ot.Exec(c, "update orders set status = $2, accrual = $3 where order_id = $1", orderId, status, accrual)
	return err
}
func (ot *ordersTable) GetOrderByOrderId(orderId string, c context.Context) (*types.Order, error) {
	r := ot.QueryRow(c, "select * from orders where order_id = $1", orderId)
	order := new(types.Order)
	err := r.Scan(&order.Id, &order.UserLogin, &order.OrderId, &order.Status, &order.Accrual, &order.UploadedAt)
	return order, err
}
func (ot *ordersTable) GetAllOrdersByLogin(login string, c context.Context) ([]types.Order, error) {
	rows, err := ot.Query(c, "select * from orders where user_login = $1", login)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[types.Order])
}

func InitDataBase() (UserTable, OrderTable) {
	configSql, err := os.ReadFile("cmd/gophermart/config/init.sql")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:pescola@localhost:5432/gophermart")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	c, cf := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cf()
	_, err = conn.Exec(c, string(configSql))
	if err != nil {
		log.Fatal(err)
	}
	return &userTable{conn}, &ordersTable{conn}
}
