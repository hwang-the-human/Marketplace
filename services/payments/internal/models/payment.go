package models

import (
	"github.com/google/uuid"
	"time"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentCompleted PaymentStatus = "completed"
	PaymentFailed    PaymentStatus = "failed"
)

type Payment struct {
	ID        uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderID   uuid.UUID     `gorm:"type:uuid" json:"order_id"` // ID заказа
	Amount    float64       `json:"amount"`                    // Сумма оплаты
	Method    string        `gorm:"size:50" json:"method"`     // Способ оплаты (например, "credit card", "paypal")
	Status    PaymentStatus `gorm:"size:50" json:"status"`     // Статус оплаты
	CreatedAt time.Time     `json:"created_at"`
}
