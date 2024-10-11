package models

import "time"

type OutboxMessage struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	EventType      string    `gorm:"size:255"`
	Payload        []byte    `gorm:"type:bytea"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Processed      bool      `gorm:"default:false"`
	IdempotencyKey *string   `gorm:"size:255;unique"`
}
