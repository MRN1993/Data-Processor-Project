package logic

import (
	"context"
	"data-processor-project/internal/logs"
	"database/sql"
	"errors"
	"time"

	"go.uber.org/zap"
)

type RequestUpdateTime struct {
    UpdatedAt time.Time
}

// RegisterRequest adds a new request to the database.
func RegisterRequest(db *sql.DB, requestID int, userID int, data string) error {
    logs.Logger.Info("Adding new request", zap.Int("requestID", requestID), zap.Int("userID", userID), zap.String("data", data))

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `INSERT INTO requests (request_id, user_id, data, received_at) VALUES (?, ?, ?, ?)`
    _, err := db.ExecContext(ctx, query, requestID, userID, data, time.Now())
    if err != nil {
        logs.Logger.Error("Failed to add request", zap.Error(err))
        return err
    }

    logs.Logger.Info("Request added successfully")
    return nil
}

// ValidateRequest validates the incoming request data.
func ValidateRequest(db *sql.DB, requestID int, userID int, data string) error {
    logs.Logger.Info("Validating request data", zap.Int("userID", userID), zap.String("data", data))

    if data == "" {
        logs.Logger.Error("Validation failed: data is empty", zap.Int("userID", userID))
        return errors.New("data cannot be empty")
    }

    // Check if user exists
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
    err := db.QueryRow(query, userID).Scan(&exists)
    if err != nil {
        logs.Logger.Error("Failed to check if user exists", zap.Error(err))
        return err
    }

    if !exists {
        logs.Logger.Error("Validation failed: user does not exist", zap.Int("userID", userID))
        return errors.New("user does not exist")
    }

    logs.Logger.Info("Request data validation successful", zap.Int("userID", userID), zap.String("data", data))
    return nil

}

