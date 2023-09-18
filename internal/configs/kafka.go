package configs

import (
	"github.com/segmentio/kafka-go"
	"os"
)

func InitKafka() *kafka.Reader {
	config := kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER_URL")},
		GroupID: "my-group",
		Topic:   "FIO",
	}

	reader := kafka.NewReader(config)
	return reader
}
