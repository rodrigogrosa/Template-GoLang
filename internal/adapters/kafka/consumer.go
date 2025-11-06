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

type kafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topic    string
	groupID  string
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(brokers []string, topic, groupID string) (ports.MessageConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	return &kafkaConsumer{
		consumer: consumer,
		topic:    topic,
		groupID:  groupID,
	}, nil
}

func (c *kafkaConsumer) Start(ctx context.Context) error {
	handler := &consumerGroupHandler{}

	go func() {
		for {
			if err := c.consumer.Consume(ctx, []string{c.topic}, handler); err != nil {
				logger.Logger.Error().Err(err).Msg("Error consuming messages")
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	logger.Logger.Info().Str("topic", c.topic).Str("group", c.groupID).Msg("Kafka consumer started")
	return nil
}

func (c *kafkaConsumer) Close() error {
	return c.consumer.Close()
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct{}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.handleMessage(message)
		session.MarkMessage(message, "")
	}
	return nil
}

func (h *consumerGroupHandler) handleMessage(msg *sarama.ConsumerMessage) {
	eventType := string(msg.Key)

	logger.Logger.Debug().
		Str("topic", msg.Topic).
		Str("event_type", eventType).
		Int32("partition", msg.Partition).
		Int64("offset", msg.Offset).
		Msg("Message received")

	switch eventType {
	case "item.created", "item.updated":
		var item domain.Item
		if err := json.Unmarshal(msg.Value, &item); err != nil {
			logger.Logger.Error().Err(err).Msg("Failed to unmarshal item")
			metrics.KafkaMessagesConsumed.WithLabelValues(msg.Topic, "error").Inc()
			return
		}
		logger.Logger.Info().Str("item_id", item.ID).Str("event", eventType).Msg("Item event processed")
	case "item.deleted":
		var data map[string]string
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logger.Logger.Error().Err(err).Msg("Failed to unmarshal delete event")
			metrics.KafkaMessagesConsumed.WithLabelValues(msg.Topic, "error").Inc()
			return
		}
		logger.Logger.Info().Str("item_id", data["id"]).Msg("Item deleted event processed")
	default:
		logger.Logger.Warn().Str("event_type", eventType).Msg("Unknown event type")
	}

	metrics.KafkaMessagesConsumed.WithLabelValues(msg.Topic, "success").Inc()
}
