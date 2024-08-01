package app

import (
	"github.com/segmentio/kafka-go"

	"github.com/crseat/example-data-pipeline/internal/domain"
)

// ProducerService provides methods for processing post data by writing it to Kafka.
// It uses a Producer instance to handle the actual communication with Kafka.
type ProducerService struct {
	producer Producer
}

// Producer defines the methods required for a Producer.
// Implementations must provide methods to write post data and raw Kafka messages to Kafka.
type Producer interface {
	WritePostDataToKafka(postData domain.PostData) error
	WriteMessageToKafka(message kafka.Message) error
}

// NewProducerService creates a new instance of ProducerService with the provided Producer.
func NewProducerService(producer Producer) *ProducerService {
	return &ProducerService{producer: producer}
}

// ProcessPostData writes the given POST data to Kafka for processing.
func (s *ProducerService) ProcessPostData(postData domain.PostData) error {
	return s.producer.WritePostDataToKafka(postData)
}

// SetProducer sets a new Producer instance for the ProducerService.
// This method is intended only for testing purposes.
func (s *ProducerService) SetProducer(producer Producer) {
	s.producer = producer
}
