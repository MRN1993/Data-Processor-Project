package my_kafka

import (
	"data-processor-project/internal/domain/logic"
	"data-processor-project/internal/logs"
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	db       *sql.DB
}

func NewKafkaConsumer(db *sql.DB, brokers, groupID string) (*KafkaConsumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: consumer, db: db}, nil
}

func (kc *KafkaConsumer) processMessage(message *kafka.Message) error {
	var requestData map[string]interface{}
	if err := json.Unmarshal(message.Value, &requestData); err != nil {
		logs.Logger.Error("Failed to unmarshal message", zap.Error(err))
		return err
	}
 
	userIDFloat, ok := requestData["userID"].(float64)
	if !ok {
		logs.Logger.Error("Invalid userID format")
		return errors.New("invalid userID format")
	}
	userID := int(userIDFloat)
 
	requestIDFloat, ok := requestData["id"].(float64)
	if !ok {
		logs.Logger.Error("Invalid requestID format")
		return errors.New("invalid requestID format")
	}
	requestID := int(requestIDFloat)
 
	data, ok := requestData["data"].(string)
	if !ok {
		logs.Logger.Error("Invalid data format")
		return errors.New("invalid data format")
	}
 
	if err := logic.RegisterRequest(kc.db, requestID, userID, data); err != nil {
		logs.Logger.Error("Failed to add request", zap.Error(err))
		return err
	}
 
	return nil
 }

func (kc *KafkaConsumer) StartConsuming(topic string) {
	defer kc.consumer.Close() 
	kc.consumer.Subscribe(topic, nil)
	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err == nil {
			if err := kc.processMessage(msg); err != nil {
				logs.Logger.Error("Failed to process message", zap.Error(err))
			}
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}