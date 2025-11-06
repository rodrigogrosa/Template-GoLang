package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/rodrigogrosa/Template-GoLang/internal/domain"
	"github.com/rodrigogrosa/Template-GoLang/pkg/logger"
)

// Publisher implements EventPublisher using Kafka
type Publisher struct {
	producer sarama.SyncProducer
	topic    string
	logger   *logger.Logger
}

// NewPublisher creates a new Kafka publisher
func NewPublisher(brokers []string, topic string, log *logger.Logger) (*Publisher, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Publisher{
		producer: producer,
		topic:    topic,
		logger:   log,
	}, nil
}

// PublishItemEvent publishes an item event to Kafka
func (p *Publisher) PublishItemEvent(ctx context.Context, event *domain.ItemEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(event.ItemID),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.logger.Error("Failed to send message to Kafka", err)
		return fmt.Errorf("failed to send message: %w", err)
	}

	p.logger.Debug("Message sent to Kafka", "partition", partition, "offset", offset)
	return nil
}

// Close closes the producer
func (p *Publisher) Close() error {
	return p.producer.Close()
}

// Consumer implements EventConsumer using Kafka
type Consumer struct {
	consumer sarama.ConsumerGroup
	topic    string
	logger   *logger.Logger
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, topic, groupID string, log *logger.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_6_0_0

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		topic:    topic,
		logger:   log,
	}, nil
}

// Subscribe subscribes to item events
func (c *Consumer) Subscribe(ctx context.Context, handler func(*domain.ItemEvent) error) error {
	h := &consumerHandler{
		handler: handler,
		logger:  c.logger,
	}

	topics := []string{c.topic}
	for {
		if err := c.consumer.Consume(ctx, topics, h); err != nil {
			c.logger.Error("Error from consumer", err)
			return err
		}

		// Check if context was cancelled
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return c.consumer.Close()
}

type consumerHandler struct {
	handler func(*domain.ItemEvent) error
	logger  *logger.Logger
}

func (h *consumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var event domain.ItemEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			h.logger.Error("Failed to unmarshal event", err)
			continue
		}

		if err := h.handler(&event); err != nil {
			h.logger.Error("Failed to handle event", err)
			continue
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

// NoOpPublisher is a no-op implementation when Kafka is disabled
type NoOpPublisher struct{}

func NewNoOpPublisher() *NoOpPublisher {
	return &NoOpPublisher{}
}

func (p *NoOpPublisher) PublishItemEvent(ctx context.Context, event *domain.ItemEvent) error {
	return nil
}

func (p *NoOpPublisher) Close() error {
	return nil
}
