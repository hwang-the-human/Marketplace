package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"reflect"
)

type Producer interface {
	Emit(topic string, message interface{}) error
	Close() error
}

type producer struct {
	Producer sarama.SyncProducer
}

func NewProducer(brokers []string) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	kafkaProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logrus.Errorf("Failed to create Kafka producer: %v", err)
		return nil, err
	}

	logrus.Info("Kafka producer created successfully")
	return &producer{Producer: kafkaProducer}, nil
}

func (p *producer) Emit(topic string, message interface{}) error {
	var messageBytes []byte

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

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(messageBytes),
	}

	partition, offset, err := p.Producer.SendMessage(msg)
	if err != nil {
		logrus.Errorf("Failed to send message to Kafka topic %s: %v", topic, err)
		return err
	}

	logrus.Infof("Message sent to Kafka topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

func (p *producer) Close() error {
	err := p.Producer.Close()
	if err != nil {
		logrus.Errorf("Failed to close Kafka producer: %v", err)
		return err
	}
	logrus.Info("Kafka producer closed successfully")
	return nil
}
