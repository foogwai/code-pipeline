package infrastructure

import (
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"

	"github.com/crseat/example-data-pipeline/internal/adapters/http"
	"github.com/crseat/example-data-pipeline/internal/adapters/kafka"
	"github.com/crseat/example-data-pipeline/internal/adapters/repositories"
	"github.com/crseat/example-data-pipeline/internal/app"
)

type (
	// CustomValidator implements the echo.Validator interface for custom validation using the go-playground/validator package.
	CustomValidator struct {
		validator *validator.Validate
	}
)

// Validate validates the given struct using the custom validator.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// StartServer initializes and starts the web server and Kafka consumer.
//
// It performs the following tasks:
//   - Loads configuration settings from environment variables.
//   - Sets up an Echo web server with middleware and custom validator.
//   - Initializes Kafka producer and Aerospike repository.
//   - Creates and registers HTTP handlers.
//   - Initializes a Kafka consumer and starts it in a separate goroutine.
//   - Starts the Echo server and waits for the Kafka consumer goroutine to finish.
func StartServer() {
	// Load configuration
	config := LoadConfig()

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register Validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Initialize Kafka producer
	producer := kafka.NewKafkaProducer(config.KafkaBrokers, config.KafkaTopic)
	defer producer.Close()

	// Initialize Aerospike repository
	repository, err := repositories.NewAerospikeRepository(config.AerospikeHost, config.AerospikePort)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Initialize service
	service := app.NewProducerService(producer)

	// Initialize handler and register routes
	handler := http.NewHandler(service)
	handler.RegisterRoutes(e)

	// Initialize Kafka consumer
	consumer := kafka.NewKafkaConsumer(config.KafkaBrokers, config.KafkaTopic, "example-consumer-group")
	dltProducer := kafka.NewKafkaProducer(config.KafkaBrokers, config.KafkaDltTopic)
	consumerService := app.NewConsumerService(consumer, dltProducer, repository)

	// Start Kafka consumer in a separate goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumerService.ConsumeMessages()
	}()

	// Start the Echo server
	e.Logger.Fatal(e.Start(config.ServerPort))

	// Wait for Kafka consumer goroutine to finish
	wg.Wait()
}
