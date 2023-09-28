package services

import (
	"TestCase/internal/configs"
	"TestCase/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log" // Импортируйте zerolog
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
			log.Error().Err(err).Msg("Error reading Kafka message")
			return
		}
		var person models.Input
		if err := json.Unmarshal(msg.Value, &person); err != nil {
			log.Error().Err(err).Msg("Error parsing Kafka message")
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
	if !isValidPerson(person) {
		log.Error().Msg("Invalid person data: missing required fields or incorrect format")
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

	if err := ks.ProduceMessage(fioFailedTopic, message); err != nil {
		log.Error().Err(err).Msg("Error producing message")
	}
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
