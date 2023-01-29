package repoitory

import (
	"context"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
	"time"
)

type withdrawTable struct {
	*pgx.Conn
}

func (wt *withdrawTable) CreateWithdraw(login, orderID string, sum float64, c context.Context) error {
	_, err := wt.Exec(c, "insert into withdraws(user_login, order_id, sum, processed_at) values ($1, $2, $3, $4)", login, orderID, sum, time.Now())
	return err
}

func (wt *withdrawTable) GetAllWithdrawalsByLogin(login string, c context.Context) ([]types.Withdraw, error) {
	rows, err := wt.Query(c, "select * from withdraws where user_login = $1", login)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[types.Withdraw])
}
