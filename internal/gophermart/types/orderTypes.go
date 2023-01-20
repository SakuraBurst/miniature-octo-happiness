package types

import "time"

type OrderStatus string
type LoyaltyServiceStatus string

const (
	NewOrder        OrderStatus = "NEW"
	ProcessingOrder OrderStatus = "PROCESSING"
	InvalidOrder    OrderStatus = "INVALID"
	ProcessedOrder  OrderStatus = "PROCESSED"
)
const (
	LoyaltyServiceRegistered LoyaltyServiceStatus = "REGISTERED"
	LoyaltyServiceInvalid    LoyaltyServiceStatus = "INVALID"
	LoyaltyServiceProcessing LoyaltyServiceStatus = "PROCESSING"
	LoyaltyServiceProcessed  LoyaltyServiceStatus = "PROCESSED"
)

type Order struct {
	Id         int         `json:"-"`
	OrderId    string      `json:"number" db:"order_id"`
	UserLogin  string      `json:"-" db:"user_login"`
	Status     OrderStatus `json:"status,omitempty"`
	Accrual    float64     `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at" db:"uploaded_at"`
}

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type Withdraw struct {
	Id          int       `json:"-"`
	UserLogin   string    `json:"-" db:"user_login"`
	Order       string    `json:"order" db:"order_id"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

type LoyaltyServiceResponse struct {
	Order   string               `json:"order"`
	Status  LoyaltyServiceStatus `json:"status"`
	Accrual float64              `json:"accrual"`
}
