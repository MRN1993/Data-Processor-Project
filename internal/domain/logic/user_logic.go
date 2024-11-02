package logic

import (
    "context"
    "database/sql"
    "errors"
    "time"
    "data-processor-project/internal/logs"
    "go.uber.org/zap"
)

type UserLimits struct {
    RequestLimitPerMinute int
    MonthlyDataLimit      int
    UsedData              int
    RequestCount          int
    LastRequestTime       time.Time
}


func RegisterUserInDB(db *sql.DB, quota, monthlyDataLimit, requestLimitPerMinute int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `INSERT INTO users (quota, monthly_data_limit, request_limit_per_minute, used_data, request_count, last_request_time) VALUES (?, ?, ?, 0, 0, ?)`
    _, err := db.ExecContext(ctx, query, quota, monthlyDataLimit, requestLimitPerMinute, time.Now())
    return err
}


func CheckUserLimits(db *sql.DB, userID, dataSize int) error {
    user, err := RetrieveUserLimits(db, userID)
    if err != nil {
        logs.Logger.Error("Failed to retrieve user data", zap.Error(err), zap.Int("userID", userID))
        return err
    }

    if user.UsedData+dataSize > user.MonthlyDataLimit {
        logs.Logger.Warn("User exceeded monthly data limit", zap.Int("userID", userID), zap.Int("usedData", user.UsedData), zap.Int("dataSize", dataSize))
        return errors.New("user has exceeded monthly data limit")
    }

    if !CheckRequestLimit(user) {
        logs.Logger.Warn("User exceeded request limit per minute", zap.Int("userID", userID), zap.Int("requestCount", user.RequestCount))
        return errors.New("user has exceeded request limit per minute")
    }

    return UpdateUserQuota(db, userID, dataSize, user)
}


func RetrieveUserLimits(db *sql.DB, userID int) (UserLimits, error) {
    var user UserLimits
    query := `SELECT request_limit_per_minute, monthly_data_limit, used_data, request_count, last_request_time FROM users WHERE id = ?`
    row := db.QueryRow(query, userID)
    err := row.Scan(&user.RequestLimitPerMinute, &user.MonthlyDataLimit, &user.UsedData, &user.RequestCount, &user.LastRequestTime)
    return user, err
}


func CheckRequestLimit(user UserLimits) bool {
    now := time.Now()
    if now.Sub(user.LastRequestTime) < time.Minute {
        return user.RequestCount+1 <= user.RequestLimitPerMinute
    }
    return true
}


func UpdateUserQuota(db *sql.DB, userID, dataSize int, user UserLimits) error {
    now := time.Now()
    newUsedData := user.UsedData + dataSize
    newRequestCount := user.RequestCount + 1

    if now.Sub(user.LastRequestTime) >= time.Minute {
        newRequestCount = 0
    }

    query := `UPDATE users SET used_data = ?, request_count = ?, last_request_time = ? WHERE id = ?`
    _, err := db.Exec(query, newUsedData, newRequestCount, now, userID)
    if err != nil {
        logs.Logger.Error("Failed to update user data", zap.Error(err), zap.Int("userID", userID))
    }
    return err
}
