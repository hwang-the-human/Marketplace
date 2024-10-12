package models

import (
	"github.com/google/uuid"
)

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid" json:"order_id"`
	ProductID uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
}
