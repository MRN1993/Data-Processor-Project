package services

import (
	"context"
	"data-processor-project/internal/domain/logic"
	"data-processor-project/internal/logs"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RequestService struct {
	db           *sql.DB
	kafkaService *KafkaService
	redisClient  *redis.Client
}

func NewRequestService(db *sql.DB, kafkaService *KafkaService, redisClient *redis.Client) *RequestService {
	return &RequestService{db: db, kafkaService: kafkaService, redisClient: redisClient}
}

func acquireLock(ctx context.Context, rdb *redis.Client, key string, expiration time.Duration) (bool, error) {
	result, err := rdb.SetNX(ctx, key, "locked", expiration).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

func releaseLock(ctx context.Context, rdb *redis.Client, key string) error {
	_, err := rdb.Del(ctx, key).Result()
	return err
}


func (s *RequestService) ProcessRequest(requestID int, userID int, data string) error {
    logs.Logger.Info("Processing request")

    ctx := context.Background()
    lockKey := "request-lock:" + strconv.Itoa(requestID)

    // Attempt to acquire lock for 1 minute
    locked, err := acquireLock(ctx, s.redisClient, lockKey, 1*time.Minute)
    if err != nil {
        logs.Logger.Error("Failed to acquire lock", zap.Error(err))
        return errors.New("failed to acquire lock")
    }

    // If lock is acquired, proceed to process the request
    if locked {
        defer releaseLock(ctx, s.redisClient, lockKey) // Release lock at the end

        // Check if the request was processed within the last minute
        lastProcessedKey := "last-processed:" + strconv.Itoa(requestID)
        lastProcessed, err := s.redisClient.Get(ctx, lastProcessedKey).Result()
        if err == nil {
            lastProcessedTime, _ := strconv.ParseInt(lastProcessed, 10, 64)
            if time.Now().Unix()-lastProcessedTime < 60 {
                logs.Logger.Warn("Duplicate request within one minute, skipping processing", zap.Int("requestID", requestID))
                return errors.New("duplicate request within one minute")
            }
        }

        // Validate request
		if err := logic.ValidateRequest(s.db,requestID, userID, data); err != nil {
			logs.Logger.Error("Validation failed", zap.Error(err))
			return err
		}


		dataSize := len(data)
		if err := logic.CheckUserLimits(s.db, userID, dataSize); err != nil {
			logs.Logger.Warn("User limit check failed", zap.Error(err))
			return err
		}

		userLimits, err := logic.RetrieveUserLimits(s.db, userID)
		if err != nil {
			logs.Logger.Error("Failed to retrieve user limits", zap.Error(err))
			return err
		}

		if err := logic.UpdateUserQuota(s.db, userID, dataSize, userLimits); err != nil {
			logs.Logger.Error("Failed to update user quota", zap.Error(err))
			return err
		}

     
        // Register request in DB and update `updated_at` timestamp
        if err = logic.RegisterRequest(s.db, requestID, userID, data); err != nil {
            logs.Logger.Error("Failed to add request", zap.Error(err))
            return err
        }

        // Update last processed time in Redis
        s.redisClient.Set(ctx, lastProcessedKey, strconv.FormatInt(time.Now().Unix(), 10), 0)

        // Prepare data for Kafka
        requestData := map[string]interface{}{
            "id":     requestID,
            "userID": userID,
            "data":   data,
        }

        requestDataJSON, _ := json.Marshal(requestData)

        go func() {
            // Send to Kafka
            err := s.kafkaService.SendRequestToKafka("unique_data_topic", string(requestDataJSON))
            if err != nil {
                logs.Logger.Error("Failed to send request to Kafka", zap.Error(err))
            }
        }()
    } else {
        logs.Logger.Warn("Could not acquire lock for request")
    }

    logs.Logger.Info("Request processing completed successfully")
    return nil
}
