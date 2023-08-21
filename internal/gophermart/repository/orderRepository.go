package repository

import (
	"context"
	"github.com/SakuraBurst/miniature-octo-happiness/internal/gophermart/types"
	"github.com/jackc/pgx/v5"
	"time"
)

type ordersTable struct {
	*pgx.Conn
}

func (ot *ordersTable) CreateOrder(login, orderID string, c context.Context) error {
	_, err := ot.Exec(c, "insert into orders(status, user_login, order_id, uploaded_at) values ( $1,$2, $3, $4)", types.NewOrder, login, orderID, time.Now())
	return err
}

func (ot *ordersTable) UpdateOrder(orderID string, status types.OrderStatus, accrual float64, c context.Context) error {
	_, err := ot.Exec(c, "update orders set status = $2, accrual = $3 where order_id = $1", orderID, status, accrual)
	return err
}
func (ot *ordersTable) GetOrderByOrderID(orderID string, c context.Context) (*types.Order, error) {
	r := ot.QueryRow(c, "select * from orders where order_id = $1", orderID)
	order := new(types.Order)
	err := r.Scan(&order.ID, &order.UserLogin, &order.OrderID, &order.Status, &order.Accrual, &order.UploadedAt)
	return order, err
}
func (ot *ordersTable) GetAllOrdersByLogin(login string, c context.Context) ([]types.Order, error) {
	rows, err := ot.Query(c, "select * from orders where user_login = $1", login)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[types.Order])
}
