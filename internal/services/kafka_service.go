package services

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaService struct {
	Consumer *kafka.Consumer
	Producer *kafka.Producer
}

func NewKafkaService(consumer *kafka.Consumer, producer *kafka.Producer) *KafkaService {
	return &KafkaService{
		Consumer: consumer,
		Producer: producer,
	}
}

func (ks *KafkaService) ConsumeMessages() {
	// Здесь обрабатывайте полученные сообщения из Kafka
	for {
		select {
		case msg := <-ks.Consumer.Events():
			switch ev := msg.(type) {
			case *kafka.Message:
				// Обработка Kafka сообщения
				fmt.Printf("Received message: %s\n", string(ev.Value))
				// Далее обрабатывайте сообщение и выполняйте необходимые действия
			case kafka.Error:
				// Обработка ошибок Kafka
				fmt.Printf("Kafka error: %v\n", ev)
			}
		}
	}
}

func (ks *KafkaService) ProduceMessage(topic string, message []byte) error {
	// Отправка сообщения в Kafka
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
