package models

import (
	"github.com/google/uuid"
	"time"
)

type Cart struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	BuyerID   uuid.UUID `gorm:"type:uuid" json:"buyer_id"`
	CreatedAt time.Time `json:"created_at"`
}
