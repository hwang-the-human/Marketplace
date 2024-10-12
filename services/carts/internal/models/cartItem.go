package models

import (
	"github.com/google/uuid"
)

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	CartID    uuid.UUID `gorm:"type:uuid" json:"cart_id"`
	ProductID uuid.UUID `gorm:"type:uuid" json:"product_id"`
	Quantity  int       `json:"quantity"`
}
