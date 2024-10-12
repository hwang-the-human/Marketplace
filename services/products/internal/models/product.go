package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"size:255" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `json:"price"`
	SellerID    uuid.UUID `gorm:"type:uuid" json:"seller_id"`
	Stock       int       `json:"stock"`
	ImageURL    string    `gorm:"size:255" json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
