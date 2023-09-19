package configs

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
)

func InitKafkaReader() *kafka.Reader {
	config := kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKER_URL")},
		GroupID: "my-group",
		Topic:   "FIO",
	}
	reader := kafka.NewReader(config)
	return reader
}

func InitKafkaWriter() *kafka.Writer {
	address := []string{os.Getenv("KAFKA_BROKER_URL")}
	topicName := "FIO"
	w := &kafka.Writer{
		Addr:     kafka.TCP(address...),
		Topic:    topicName,
		Balancer: &kafka.LeastBytes{},
	}
	err := createKafkaTopic(topicName, address[0])
	if err != nil {
		log.Fatalf("Error creating Kafka topic: %v", err)
	}
	message := []byte(`{
        "name": "Dmitriy",
        "surname": "Ushakov",
        "patronymic": "Vasilevich"
    }`)
	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   nil,
			Value: message,
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	} else {
		fmt.Println("Message sent to Kafka successfully")
	}
	return w
}

func createKafkaTopic(topicName, brokerAddress string) error {
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	config := kafka.TopicConfig{
		Topic:             topicName,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	err = conn.CreateTopics(config)
	if err != nil {
		return err
	}

	fmt.Printf("Kafka topic '%s' created successfully\n", topicName)
	return nil
}
