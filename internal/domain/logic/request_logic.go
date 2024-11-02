package logic

import (
	"context"
	"data-processor-project/internal/domain/models"
	"data-processor-project/internal/logs"
	"database/sql"
	"errors"
	"time"

	"go.uber.org/zap"
)

// AddRequest adds a new request to the database.
func RegisterRequest(db *sql.DB, UserID int,Data string) error {
    logs.Logger.Info("Adding new request", zap.Int("userID", UserID), zap.String("data", Data))

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    query := `INSERT INTO requests (user_id, data,received_at) VALUES (?, ?, ?)`
    _, err := db.ExecContext(ctx, query, UserID, Data, time.Now())
    if err != nil {
        logs.Logger.Error("Failed to add request", zap.Error(err))
        return err
    }

    logs.Logger.Info("Request added successfully")
    return nil
}

// ValidateRequest validates the incoming request data.
func ValidateRequest(db *sql.DB, userID int, data string) error {
    logs.Logger.Info("Validating request data", zap.Int("userID", userID), zap.String("data", data))

    if data == "" {
        logs.Logger.Error("Validation failed: data is empty", zap.Int("userID", userID))
        return errors.New("data cannot be empty")
    }

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

// CheckDuplicate checks if the request is a duplicate.
func CheckDuplicate(db *sql.DB, userID int, data string) (bool, error) {
    logs.Logger.Info("Checking for duplicate request", zap.Int("userID", userID), zap.String("data", data))

    var count int
    query := `SELECT COUNT(*) FROM requests WHERE user_id = ? AND data = ?`
    err := db.QueryRow(query, userID, data).Scan(&count)
    if err != nil {
        logs.Logger.Error("Failed to check for duplicate request", zap.Error(err))
        return false, err
    }

    if count > 0 {
        logs.Logger.Warn("Duplicate request detected", zap.Int("userID", userID), zap.String("data", data))
        return true, nil
    }

    logs.Logger.Info("No duplicate request found", zap.Int("userID", userID), zap.String("data", data))
    return false, nil
}

// CheckUserQuota checks if the user has exceeded their request quota.
func CheckUserQuota(db *sql.DB, userID int) error {
    logs.Logger.Info("Checking user quota", zap.Int("userID", userID))

    var user models.User // Assuming you have a User struct
    query := `SELECT quota FROM users WHERE id = ?`
    err := db.QueryRow(query, userID).Scan(&user.Quota)
    if err != nil {
        logs.Logger.Error("Failed to retrieve user quota", zap.Error(err))
        return err
    }

    var requestCount int
    query = `SELECT COUNT(*) FROM requests WHERE user_id = ?`
    err = db.QueryRow(query, userID).Scan(&requestCount)
    if err != nil {
        logs.Logger.Error("Failed to count user requests", zap.Error(err))
        return err
    }

    if requestCount >= user.Quota {
        logs.Logger.Warn("User has exceeded request quota", zap.Int("userID", userID), zap.Int("requestCount", requestCount), zap.Int("quota", user.Quota))
        return errors.New("user has exceeded request quota")
    }

    logs.Logger.Info("User quota check passed", zap.Int("userID", userID), zap.Int("requestCount", requestCount), zap.Int("quota", user.Quota))
    return nil
}
