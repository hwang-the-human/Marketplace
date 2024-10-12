package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Consumer interface {
	Consume(ctx context.Context, topics []string, handler func(message *sarama.ConsumerMessage)) error
	Close() error
}

type consumer struct {
	ConsumerGroup sarama.ConsumerGroup
}

func NewConsumer(brokers []string, groupID string) (Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRange(),
	}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	kafkaConsumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		logrus.Errorf("Failed to create Kafka consumer group: %v", err)
		return nil, err
	}

	logrus.Info("Kafka consumer group created successfully")
	return &consumer{ConsumerGroup: kafkaConsumer}, nil
}

func (c *consumer) Consume(ctx context.Context, topics []string, handler func(message *sarama.ConsumerMessage)) error {
	for {
		err := c.ConsumerGroup.Consume(ctx, topics, &consumerGroupHandler{handler: handler})
		if err != nil {
			logrus.Errorf("Error during consumption: %v", err)
			return err
		}
		if ctx.Err() != nil {
			logrus.Infof("Context canceled: %v", ctx.Err())
			return ctx.Err()
		}
	}
}

func (c *consumer) Close() error {
	err := c.ConsumerGroup.Close()
	if err != nil {
		logrus.Errorf("Failed to close Kafka consumer: %v", err)
		return err
	}
	logrus.Info("Kafka consumer closed successfully")
	return nil
}

type consumerGroupHandler struct {
	handler func(message *sarama.ConsumerMessage)
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		logrus.Infof("Message consumed from topic %s, partition %d, offset %d", message.Topic, message.Partition, message.Offset)
		h.handler(message)
		session.MarkMessage(message, "")
	}
	return nil
}
