package configs

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
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
		log.Fatal().Err(err).Msg("Error creating Kafka topic")
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
		log.Fatal().Err(err).Msg("Failed to write messages")
	} else {
		log.Info().Msg("Message sent to Kafka successfully")
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

	log.Info().Msgf("Kafka topic '%s' created successfully", topicName)
	return nil
}
