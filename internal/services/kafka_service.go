package services

import (
	"TestCase/internal/configs"
	"TestCase/internal/models"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaService struct {
	Consumer          *kafka.Consumer
	Producer          *kafka.Producer
	enrichmentService *EnrichmentService
}

func NewKafkaService(consumer *kafka.Consumer, producer *kafka.Producer) *KafkaService {
	return &KafkaService{
		Consumer: consumer,
		Producer: producer,
	}
}

func (ks *KafkaService) ConsumeMessages() {
	consumer := configs.InitKafka()
	defer consumer.Close()
	for {
		select {
		case msg := <-consumer.Events():
			switch ev := msg.(type) {
			case *kafka.Message:
				go ks.ProcessFIOMessage(ev.Value)
			case kafka.Error:
				fmt.Printf("Kafka error: %v\n", ev)
			}
		}
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
	deliveryChan := make(chan kafka.Event)

	err := ks.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)

	if err != nil {
		return err
	}
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}
	return nil
}
