package models

import "time"

type OutboxMessage struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	EventType      string    `gorm:"size:255"`
	Payload        string    `gorm:"type:jsonb"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Processed      bool      `gorm:"default:false"`
	IdempotencyKey *string   `gorm:"size:255;unique"`
}
