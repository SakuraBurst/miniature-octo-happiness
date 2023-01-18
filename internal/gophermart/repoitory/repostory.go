package repoitory

import (
	"context"
	"fmt"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
	"os"
	"time"
)

type userTable struct {
	*pgx.Conn
}

type UserTable interface {
	CreateUser(login, hashedPassword string, c context.Context) error
	GetUser(login string, c context.Context) (*types.User, error)
	UpdateBalance(login string, newBalance int, c context.Context) error
}

func (ut *userTable) CreateUser(login, hashedPassword string, c context.Context) error {
	_, err := ut.Exec(c, "insert into users(login, password) values ($1, $2)", login, hashedPassword)
	return err
}

func (ut *userTable) GetUser(login string, c context.Context) (*types.User, error) {
	r := ut.QueryRow(c, "select * from users where login = $1", login)
	user := new(types.User)
	err := r.Scan(&user.Id, &user.Login, &user.Password, &user.Balance)
	return user, err
}

func (ut *userTable) UpdateBalance(login string, newBalance int, c context.Context) error {
	_, err := ut.Exec(c, "update users set balance = $2 where login = $1", login, newBalance)
	return err
}

type ordersTable struct {
	*pgx.Conn
}

type OrderTable interface {
	CreateOrder(login, orderId string, c context.Context) error
	UpdateOrder(orderId string, status string, accrual int, c context.Context) error
	GetAllOrdersByLogin(login string, c context.Context) ([]types.Order, error)
}

func (ot *ordersTable) CreateOrder(login, orderId string, c context.Context) error {
	_, err := ot.Exec(c, "insert into orders(status, user_login, order_id, accrual, uploaded_at) values ( 'NEW',$1, $2)", login, orderId, 0, time.Now())
	return err
}

func (ot *ordersTable) UpdateOrder(orderId string, status string, accrual int, c context.Context) error {
	_, err := ot.Exec(c, "update orders set status = $2, accrual = $3 where order_id = $1", orderId, status, accrual)
	return err
}

func (ot *ordersTable) GetAllOrdersByLogin(login string, c context.Context) ([]types.Order, error) {
	rows, err := ot.Query(c, "select * from orders where user_login = $1", login)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[types.Order])
}

func InitDataBase() (UserTable, OrderTable) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:pescola@localhost:5432/gophermart")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &userTable{conn}, &ordersTable{conn}
}
