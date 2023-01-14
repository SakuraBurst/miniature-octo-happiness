package repoitory

import (
	"context"
	"fmt"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
	"os"
	"time"
)

type DataBase struct {
	*pgx.Conn
}

type userTable struct {
	*pgx.Conn
}

type UserTable interface {
	CreateUser(login, hashedPassword string, c context.Context) error
	GetUser(login string, c context.Context) (*types.User, error)
	UpdateBalance(id, newBalance int, c context.Context) error
}

func (ut *userTable) CreateUser(login, hashedPassword string, c context.Context) error {
	_, err := ut.Exec(c, "insert into users(login, password) values ($1, $2)", login, hashedPassword)
	return err
}

func (ut *userTable) GetUser(login, c context.Context) (*types.User, error) {
	r := ut.QueryRow(c, "select 1 from users where login = $1", login)
	user := new(types.User)
	err := r.Scan(&user.Id, &user.Login, &user.Password, &user.Balance)
	return user, err
}

func (ut *userTable) UpdateBalance(id, newBalance int, c context.Context) error {
	_, err := ut.Exec(c, "update users set balance = $2 where id = $1", id, newBalance)
	return err
}

type ordersTable struct {
	*pgx.Conn
}

func (ot *ordersTable) CreateOrder(userId int, c context.Context) {
	ot.Exec(c, "insert into orders(user_id, status, uploaded_at) values ($1, 'NEW', $2)", userId, time.Now())
}

func (ot *ordersTable) UpdateOrder(id int, status string, accrual int, c context.Context) {
	ot.Exec(c, "update orders set status = $2, accrual = $3 where id = $1", id, status, accrual)
}

func InitDataBase() *DataBase {
	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://postgres:pescola@localhost:5432/gophermart"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &DataBase{conn}
}
