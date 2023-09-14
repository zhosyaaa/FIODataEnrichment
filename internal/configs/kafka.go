package configs

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func InitKafka() *kafka.Consumer {
	config := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKER_URL"), // Пример: localhost:9092
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		fmt.Printf("Error creating Kafka consumer: %v\n", err)
		return nil
	}

	topics := []string{"FIO"} // Названия топиков Kafka для подписки
	consumer.SubscribeTopics(topics, nil)

	return consumer
}
