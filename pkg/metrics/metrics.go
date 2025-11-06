package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestsTotal tracks total HTTP requests
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestDuration tracks HTTP request duration
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// ItemsCreated tracks total items created
	ItemsCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "items_created_total",
			Help: "Total number of items created",
		},
	)

	// ItemsDeleted tracks total items deleted
	ItemsDeleted = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "items_deleted_total",
			Help: "Total number of items deleted",
		},
	)

	// ItemsUpdated tracks total items updated
	ItemsUpdated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "items_updated_total",
			Help: "Total number of items updated",
		},
	)

	// KafkaMessagesPublished tracks Kafka messages published
	KafkaMessagesPublished = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_published_total",
			Help: "Total number of Kafka messages published",
		},
		[]string{"topic", "status"},
	)

	// KafkaMessagesConsumed tracks Kafka messages consumed
	KafkaMessagesConsumed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_consumed_total",
			Help: "Total number of Kafka messages consumed",
		},
		[]string{"topic", "status"},
	)
)
