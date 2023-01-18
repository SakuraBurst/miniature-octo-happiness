package types

import "time"

type OrderStatus string

const (
	New        OrderStatus = "NEW"
	Processing OrderStatus = "PROCESSING"
	Invalid    OrderStatus = "INVALID"
	Processed  OrderStatus = "PROCESSED"
)

type Order struct {
	Id         int         `json:"-"`
	OrderId    string      `json:"number" db:"order_id"`
	UserLogin  string      `json:"-" db:"user_login"`
	Status     OrderStatus `json:"status,omitempty"`
	Accrual    int         `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at" db:"uploaded_at"`
}

type WithdrawRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

type Withdraw struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
