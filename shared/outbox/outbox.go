package outbox

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"marketplace/shared/db"
	"marketplace/shared/kafka"
	"marketplace/shared/models"
	"reflect"
	"sync"
)

type Outbox interface {
	SaveMessageToOutbox(topic string, message interface{}, enableIdempotencyKey ...bool) error
	ProcessOutboxMessages() error
}

type outbox struct {
	Database db.Database
	Producer kafka.Producer
}

func NewOutbox(database db.Database, producer kafka.Producer) Outbox {
	return &outbox{Database: database, Producer: producer}
}

func (o *outbox) SaveMessageToOutbox(topic string, message interface{}, enableIdempotencyKey ...bool) error {
	var idempotencyKey *string
	var messageBytes []byte

	if len(enableIdempotencyKey) > 0 && enableIdempotencyKey[0] {
		key := uuid.New().String()
		idempotencyKey = &key
	}

	switch v := message.(type) {
	case []byte:
		messageBytes = v
	case string:
		messageBytes = []byte(v)
	default:
		var err error
		messageBytes, err = json.Marshal(message)
		if err != nil {
			logrus.Errorf("Failed to marshal message of type %s to JSON: %v", reflect.TypeOf(message), err)
			return err
		}
	}

	outboxMessage := models.OutboxMessage{
		EventType:      topic,
		Payload:        messageBytes,
		Processed:      false,
		IdempotencyKey: idempotencyKey,
	}

	if err := o.Database.GetDB().Create(&outboxMessage).Error; err != nil {
		logrus.Errorf("Failed to create outbox message: %v", err)
		return err
	}

	logrus.Infof("Outbox message created for topic %s", topic)

	return nil
}

func (o *outbox) ProcessOutboxMessages() error {
	chunkSize := 100
	workerCount := 10
	database := o.Database.GetDB()

	var totalMessages int64
	database.Model(&models.OutboxMessage{}).Where("processed = ?", false).Count(&totalMessages)

	for offset := 0; offset < int(totalMessages); offset += chunkSize {
		var messages []models.OutboxMessage
		database.Where("processed = ?", false).Offset(offset).Limit(chunkSize).Find(&messages)

		if len(messages) == 0 {
			break
		}

		messageChan := make(chan models.OutboxMessage)
		var wg sync.WaitGroup

		for i := 0; i < workerCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for msg := range messageChan {
					if o.alreadyProcessed(msg.IdempotencyKey) {
						continue
					}
					err := o.Producer.Emit(msg.EventType, msg.Payload)
					if err != nil {
						logrus.Errorf("Failed to send message to Kafka: %v", err)
						continue
					}
					database.Model(&msg).Update("processed", true)
				}
			}()
		}

		for _, message := range messages {
			messageChan <- message
		}

		close(messageChan)
		wg.Wait()
	}

	return nil
}

func (o *outbox) alreadyProcessed(idempotencyKey *string) bool {
	if idempotencyKey == nil {
		return false
	}

	var count int64
	o.Database.GetDB().Model(&models.OutboxMessage{}).Where("idempotency_key = ? AND processed = ?", idempotencyKey, true).Count(&count)
	return count > 0
}
