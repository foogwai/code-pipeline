package kafka

import (
	"github.com/segmentio/kafka-go"
)

// NewKafkaConsumer creates a new Kafka reader with the specified brokers, topic, and group ID.
// It returns a kafka.Reader configured to read messages from the given Kafka topic and consumer group.
func NewKafkaConsumer(brokers []string, topic string, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
}
