package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"marketplace/shared/models"
)

var Producer sarama.SyncProducer

func InitKafkaProducer(brokers []string) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logrus.Fatalf("Failed to start Kafka producer: %v", err)
	}

	Producer = producer
	logrus.Info("Successfully connected to Kafka")
}

func InitKafkaConsumer(brokers []string, topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		logrus.Fatalf("Failed to start Kafka consumer: %v", err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		logrus.Fatalf("Failed to start partition consumer: %v", err)
	}

	logrus.Info("Kafka consumer initialized")
	return partitionConsumer, nil
}

func SendMessage(db *gorm.DB, topic, message string, enableIdempotencyKey ...bool) error {
	var idempotencyKey *string

	if len(enableIdempotencyKey) > 0 && enableIdempotencyKey[0] {
		key := uuid.New().String()
		idempotencyKey = &key
	}

	outboxMessage := models.OutboxMessage{
		EventType:      topic,
		Payload:        message,
		Processed:      false,
		IdempotencyKey: idempotencyKey,
	}

	if err := db.Create(&outboxMessage).Error; err != nil {
		logrus.Errorf("Failed to create outbox message: %v", err)
		return err
	}

	logrus.Infof("Outbox message created for topic %s", topic)

	return nil
}

func ProcessOutboxMessages(db *gorm.DB, producer sarama.SyncProducer, params ...int) {
	chunkSize := 100
	workerCount := 10

	if len(params) > 0 {
		chunkSize = params[0]
	}
	if len(params) > 1 {
		workerCount = params[1]
	}

	var totalMessages int64
	db.Model(&models.OutboxMessage{}).Where("processed = ?", false).Count(&totalMessages)

	for offset := 0; offset < int(totalMessages); offset += chunkSize {
		var messages []models.OutboxMessage
		db.Where("processed = ?", false).Offset(offset).Limit(chunkSize).Find(&messages)

		if len(messages) == 0 {
			break
		}

		messageChan := make(chan models.OutboxMessage)

		for i := 0; i < workerCount; i++ {
			go func() {
				for msg := range messageChan {
					err := sendMessageToKafka(db, producer, msg)
					if err != nil {
						logrus.Errorf("Failed to send message to Kafka: %v", err)
						continue
					}
					db.Model(&msg).Update("processed", true)
				}
			}()
		}

		for _, message := range messages {
			messageChan <- message
		}

		close(messageChan)
	}
}

func CloseKafkaProducer() {
	if err := Producer.Close(); err != nil {
		logrus.Warnf("Failed to close Kafka producer: %v", err)
	} else {
		logrus.Info("Kafka producer closed")
	}
}

func sendMessageToKafka(db *gorm.DB, producer sarama.SyncProducer, message models.OutboxMessage) error {
	if message.IdempotencyKey != nil && alreadyProcessed(db, message.IdempotencyKey) {
		return fmt.Errorf("message with idempotency key %s already processed", *message.IdempotencyKey)
	}

	msg := &sarama.ProducerMessage{
		Topic: message.EventType,
		Value: sarama.StringEncoder(message.Payload),
	}

	_, _, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	logrus.Infof("Message sent to Kafka topic %s: %s", message.EventType, message.Payload)
	return nil
}

func alreadyProcessed(db *gorm.DB, idempotencyKey *string) bool {
	if idempotencyKey == nil {
		return false
	}

	var count int64
	db.Model(&models.OutboxMessage{}).Where("idempotency_key = ? AND processed = ?", idempotencyKey, true).Count(&count)
	return count > 0
}
