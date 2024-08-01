package infrastructure

import (
	"os"
	"strings"
)

// Config holds configuration settings for the application.
// It includes server settings, Kafka and Aerospike configurations, and application environment details.
//
// Fields:
//   - ServerPort: The port on which the server will listen for incoming connections.
//   - KafkaBrokers: A slice of Kafka broker addresses used for connecting to the Kafka cluster.
//   - KafkaTopic: The Kafka topic to which messages will be published.
//   - KafkaDltTopic: The Kafka topic used for the dead letter topic.
//   - AppEnvironment: The environment in which the application is running (e.g., development, production).
//   - AerospikeHost: The host address of the Aerospike database.
//   - AerospikePort: The port number of the Aerospike database.
type Config struct {
	ServerPort     string
	KafkaBrokers   []string
	KafkaTopic     string
	KafkaDltTopic  string
	AppEnvironment string
	AerospikeHost  string
	AerospikePort  int
}

// LoadConfig loads the configuration from environment variables and returns a Config instance.
func LoadConfig() *Config {
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", ":8080"),
		KafkaBrokers:   strings.Split(getEnv("KAFKA_BROKER", "localhost:9092"), ","),
		KafkaTopic:     getEnv("KAFKA_TOPIC", "data-pipeline-topic"),
		KafkaDltTopic:  getEnv("KAFKA_DLT_TOPIC", "data-pipeline-dlt-topic"),
		AppEnvironment: getEnv("APP_ENV", "development"),
		AerospikeHost:  getEnv("AEROSPIKE_HOST", "localhost"),
		AerospikePort:  3000,
	}
}

// getEnv retrieves the value of the environment variable named by key.
// If the environment variable is not set, it returns the defaultValue.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
