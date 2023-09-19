package services

import (
	"TestCase/internal/configs"
	"TestCase/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"sync"
)

type KafkaService struct {
	Reader            *kafka.Reader
	Writer            *kafka.Writer
	enrichmentService *EnrichmentService
}

func NewKafkaService(reader *kafka.Reader, writer *kafka.Writer, enrichmentService *EnrichmentService) *KafkaService {
	return &KafkaService{Reader: reader, Writer: writer, enrichmentService: enrichmentService}
}

func (ks *KafkaService) ConsumeMessages() {
	reader := configs.InitKafkaReader()
	defer reader.Close()
	var mu sync.Mutex

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Error reading Kafka message: %v\n", err)
			return
		}
		var person models.Input
		if err := json.Unmarshal(msg.Value, &person); err != nil {
			ks.SendToFIOfailed(fmt.Sprintf("Error parsing Kafka message: %v", err))
			return
		}
		if isValidPerson(person) {
			mu.Lock()
			ks.ProcessFIOMessage(person)
			mu.Unlock()
		}
	}
}

func (ks *KafkaService) ProcessFIOMessage(person models.Input) {
	fmt.Println("ProcessFIOMessage: ", person.Name)
	if !isValidPerson(person) {
		ks.SendToFIOfailed("Invalid person data: missing required fields or incorrect format")
		return
	}
	fio := fmt.Sprintf("%s %s %s", person.Name, person.Surname, person.Patronymic)
	ks.enrichmentService.FIOChannel <- fio
}

func isValidPerson(person models.Input) bool {
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
