package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/rodrigogrosa/Template-GoLang/internal/core/domain"
	"github.com/rodrigogrosa/Template-GoLang/internal/core/ports"
	"github.com/rodrigogrosa/Template-GoLang/pkg/logger"
	"github.com/rodrigogrosa/Template-GoLang/pkg/metrics"
)

type kafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string, topic string) (ports.MessageProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &kafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *kafkaProducer) PublishItemCreated(ctx context.Context, item *domain.Item) error {
	return p.publish("item.created", item)
}

func (p *kafkaProducer) PublishItemUpdated(ctx context.Context, item *domain.Item) error {
	return p.publish("item.updated", item)
}

func (p *kafkaProducer) PublishItemDeleted(ctx context.Context, itemID string) error {
	return p.publish("item.deleted", map[string]string{"id": itemID})
}

func (p *kafkaProducer) publish(eventType string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		metrics.KafkaMessagesPublished.WithLabelValues(p.topic, "error").Inc()
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(eventType),
		Value: sarama.ByteEncoder(payload),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		metrics.KafkaMessagesPublished.WithLabelValues(p.topic, "error").Inc()
		logger.Logger.Error().Err(err).Str("event_type", eventType).Msg("Failed to publish message")
		return fmt.Errorf("failed to send message: %w", err)
	}

	metrics.KafkaMessagesPublished.WithLabelValues(p.topic, "success").Inc()
	logger.Logger.Debug().Str("event_type", eventType).Msg("Message published")
	return nil
}

func (p *kafkaProducer) Close() error {
	return p.producer.Close()
}
