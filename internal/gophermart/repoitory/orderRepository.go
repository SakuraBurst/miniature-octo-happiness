package repoitory

import (
	"context"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
	"time"
)

type ordersTable struct {
	*pgx.Conn
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
