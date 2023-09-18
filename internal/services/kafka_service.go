package services

import (
	"TestCase/internal/configs"
	"TestCase/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type KafkaService struct {
	Reader            *kafka.Reader
	Writer            *kafka.Writer
	enrichmentService *EnrichmentService
}

func NewKafkaService(reader *kafka.Reader, writer *kafka.Writer) *KafkaService {
	return &KafkaService{
		Reader: reader,
		Writer: writer,
	}
}

func (ks *KafkaService) ConsumeMessages() {
	reader := configs.InitKafka()
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Error reading Kafka message: %v\n", err)
			return
		}

		go ks.ProcessFIOMessage(msg.Value)
	}
}

func (ks *KafkaService) ProcessFIOMessage(message []byte) {
	var person models.Person
	if err := json.Unmarshal(message, &person); err != nil {
		ks.SendToFIOfailed(fmt.Sprintf("Error parsing Kafka message: %v", err))
		return
	}

	if !isValidPerson(person) {
		ks.SendToFIOfailed("Invalid person data: missing required fields or incorrect format")
		return
	}

	ks.enrichmentService.FIOChannel <- person.Name
}

func isValidPerson(person models.Person) bool {
	if person.Name == "" || person.Surname == "" {
		return false
	}
	return true
}

func (ks *KafkaService) SendToFIOfailed(reason string) {
	fioFailedTopic := "FIO_FAILED"
	message := []byte(reason)

	ks.ProduceMessage(fioFailedTopic, message)
}

func (ks *KafkaService) ProduceMessage(topic string, message []byte) error {
	writer := ks.Writer

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   nil,
			Value: message,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
