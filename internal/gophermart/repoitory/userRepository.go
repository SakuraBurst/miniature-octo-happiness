package repoitory

import (
	"context"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
)

type userTable struct {
	*pgx.Conn
}

func (ut *userTable) CreateUser(login, hashedPassword string, c context.Context) error {
	_, err := ut.Exec(c, "insert into users(login, password) values ($1, $2)", login, hashedPassword)
	return err
}

func (ut *userTable) GetUser(login string, c context.Context) (*types.User, error) {
	r := ut.QueryRow(c, "select * from users where login = $1", login)
	user := new(types.User)
	err := r.Scan(&user.ID, &user.Login, &user.Password, &user.Balance, &user.Withdraw)
	return user, err
}

func (ut *userTable) UpdateBalanceAndWithdraw(login string, newBalance, newWithdraw float64, c context.Context) error {
	_, err := ut.Exec(c, "update users set balance = $2, withdraw = $3 where login = $1", login, newBalance, newWithdraw)
	return err
}
