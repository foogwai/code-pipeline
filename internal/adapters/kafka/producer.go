package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/crseat/example-data-pipeline/internal/domain"
)

// Producer represents a Kafka producer that writes messages to a Kafka topic.
type Producer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new Producer with the specified brokers and topic.
// The kafka writer uses the LeastBytes balancer to distribute messages and produces in batches to increase efficiency.
func NewKafkaProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    100,
			BatchTimeout: 500 * time.Millisecond,
		},
	}
}

// WritePostDataToKafka writes the given post data to the Kafka topic.
// It marshals the post data to JSON before writing it to the topic.
// Returns an error if the post data could not be marshaled or if the write operation fails.
func (p *Producer) WritePostDataToKafka(postData domain.PostData) error {
	message, err := json.Marshal(postData)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Value: message,
	})
}

// WriteMessageToKafka writes the given Kafka message to the Kafka topic.
// Returns an error if the write operation fails.
func (p *Producer) WriteMessageToKafka(message kafka.Message) error {
	return p.writer.WriteMessages(context.Background(), message)
}

// Close closes the Kafka producer, releasing any resources.
// Logs an error message if the close operation fails.
func (p *Producer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Printf("failed to close writer: %v", err)
	}
}
