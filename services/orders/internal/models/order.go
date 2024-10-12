package models

import (
	"github.com/google/uuid"
	"time"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusCompleted OrderStatus = "completed"
	StatusCanceled  OrderStatus = "canceled"
	StatusShipped   OrderStatus = "shipped"
)

type Order struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	BuyerID       uuid.UUID   `gorm:"type:uuid" json:"buyer_id"`
	Status        OrderStatus `gorm:"size:100" json:"status"`
	TotalAmount   float64     `json:"total_amount"`
	PaymentMethod string      `gorm:"size:50" json:"payment_method"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}
