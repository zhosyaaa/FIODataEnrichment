package configs

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"os"
)

func InitKafka() *kafka.Consumer {
	config := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKER_URL"),
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		//fmt.Printf("Error creating Kafka consumer: %v\n", err)
		//return nil
		panic(err)
	}

	topics := []string{"FIO"}
	consumer.SubscribeTopics(topics, nil)

	return consumer
}
