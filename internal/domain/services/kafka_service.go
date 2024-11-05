package services

import (
	"data-processor-project/internal/kafka"
	"go.uber.org/zap"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"data-processor-project/internal/logs"
)

type KafkaService struct {
	producer *kafka.Producer
}

func NewKafkaService() (*KafkaService, error) {
	producer, err := my_kafka.NewKafkaProducer("localhost:9092")
	if err != nil {
		logs.Logger.Error("Failed to create Kafka producer", zap.Error(err))
		return nil, err
	}

	logs.Logger.Info("Kafka producer created successfully")
	return &KafkaService{producer: producer}, nil
}

func (k *KafkaService) SendRequestToKafka(topic, message string) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}

	logs.Logger.Info("Sending message to Kafka", zap.String("topic", topic), zap.String("message", message))

	err := k.producer.Produce(msg, nil) 
	if err != nil {
		logs.Logger.Error("Failed to send message to Kafka", zap.Error(err), zap.String("topic", topic))
		return err
	}


	k.producer.Flush(5 * 1000)

	logs.Logger.Info("Message delivered successfully", zap.String("topic", topic))
	return nil
}

func (k *KafkaService) Close() {
	k.producer.Close()
	logs.Logger.Info("Kafka producer closed successfully")
}