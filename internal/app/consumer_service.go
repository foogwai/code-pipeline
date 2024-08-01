package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/crseat/example-data-pipeline/internal/domain"
)

// ConsumerService represents a service for consuming messages from Kafka and processing them.
type ConsumerService struct {
	reader     *kafka.Reader
	producer   Producer
	repository domain.Repository
}

// NewConsumerService creates a new instance of ConsumerService with the provided Kafka reader, producer, and repository.
func NewConsumerService(reader *kafka.Reader, producer Producer, repository domain.Repository) *ConsumerService {
	return &ConsumerService{reader: reader, repository: repository}
}

// ConsumeMessages starts consuming messages from the Kafka topic.
// It reads messages, unmarshals them into post data, and saves the post data to the repository.
// If an error occurs during saving, it writes the message to a dead letter topic.
func (s *ConsumerService) ConsumeMessages() {
	for {
		message, err := s.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}

		var postData domain.PostData
		if err := json.Unmarshal(message.Value, &postData); err != nil {
			log.Printf("error unmarshalling message: %v", err)
			continue
		}

		if err := s.repository.SavePostData(postData); err != nil {
			log.Printf("error saving POST data: %v", err)

			// Looks like our backend is down. Let's write the messages to a dead letter topic so we can re-drive them
			// to the backend after the outage.
			if dltErr := s.producer.WriteMessageToKafka(message); dltErr != nil {
				log.Printf("error writing message to Dead Letter Topic: %v", dltErr)
			}
		}
	}
}
