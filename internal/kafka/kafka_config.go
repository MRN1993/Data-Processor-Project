
package my_kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func NewKafkaProducer(brokers string) (*kafka.Producer, error) {
	config := &kafka.ConfigMap{"bootstrap.servers": brokers}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Printf("Failed to create Kafka producer: %v", err)
		return nil, err
	}
	return producer, nil
}
