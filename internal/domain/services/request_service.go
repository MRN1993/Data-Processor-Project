package services

import (
    "fmt"
    "time"
    "data-processor-project/internal/domain/models"
    "data-processor-project/internal/generator"
    "data-processor-project/internal/domain/logic" // Importing the logic package
    "go.uber.org/zap"
)

// RequestService provides functionality to manage requests.
type RequestService struct {
    requests map[string]models.Request // Store requests
    users    map[string]models.User     // Store users
    logger   *zap.Logger                // Logger for RequestService
}

// NewRequestService creates a new RequestService.
func NewRequestService(logger *zap.Logger) *RequestService {
    // Initialize users with example data
    users := map[string]models.User{
        "user1": {ID: "user1", Quota: 5}, // Example user with a quota of 5 requests
        // You can add more users here
    }

    return &RequestService{
        requests: make(map[string]models.Request),
        users:    users,
        logger:   logger,
    }
}

// HandleRequest processes a new request.
func (s *RequestService) HandleRequest(userID string, data string) error {
    // 1. Validate the request
    if err := logic.ValidateRequest(userID, data); err != nil {
        s.logger.Warn("Request validation failed", zap.String("userID", userID), zap.String("data", data), zap.Error(err))
        return err
    }

    // 2. Check for duplicates
    if logic.CheckDuplicate(s.requests, userID, data) {
        s.logger.Warn("Duplicate request detected", zap.String("userID", userID), zap.String("data", data))
        return fmt.Errorf("duplicate request detected")
    }

    // 3. Check user quota
    if err := logic.CheckUserQuota(s.users, s.requests, userID); err != nil {
        s.logger.Warn("User quota check failed", zap.String("userID", userID), zap.Error(err))
        return err
    }

    // 4. Create the request
    request := models.Request{
        ID:         generator.GenerateUUID(),
        UserID:     userID,
        Data:       data,
        ReceivedAt: time.Now(),
    }

    // Store the request in memory
    s.requests[request.ID] = request

    // Log the request processing
    s.logger.Info("Processing request", zap.String("requestID", request.ID), zap.String("userID", userID), zap.String("data", data))

    return nil
}