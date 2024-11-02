package services

import (
    "context"
    "database/sql"
    "errors"
    "time"
    "data-processor-project/internal/logs"
    "go.uber.org/zap"
)

type UserService struct {
    db     *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
    return &UserService{
        db:     db,
    }
}

// متد RegisterUser برای ثبت یک کاربر جدید
func (us *UserService) RegisterUser(quota int, monthlydatalimit int,requestlimitperminute int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    logs.Logger.Info("Attempting to register user", zap.Int("quota", quota))

    query := `INSERT INTO users (quota, monthly_data_limit, request_limit_per_minute, used_data, request_count, last_request_time) VALUES (?, ?, ?, 0, 0, ?)`
    _, err := us.db.ExecContext(ctx, query, quota, monthlydatalimit, requestlimitperminute, 0, 0, time.Now())
    if err != nil {
        logs.Logger.Error("Failed to register user", zap.Error(err))
        return errors.New("failed to register user: " + err.Error())
    }

    logs.Logger.Info("User registered successfully")
    return nil
}

// متد CheckQuota برای بررسی محدودیت‌های کاربر
func (us *UserService) CheckQuota(userID int, dataSize int) error {
    logs.Logger.Info("Checking user quota", zap.Int("userID", userID), zap.Int("dataSize", dataSize))

    var (
        requestLimitPerMinute, monthlyDataLimit, usedData, requestCount int
        lastRequestTime time.Time
    )

    query := `SELECT request_limit_per_minute, monthly_data_limit, used_data, request_count, last_request_time FROM users WHERE id = ?`
    row := us.db.QueryRow(query, userID)
    err := row.Scan(&requestLimitPerMinute, &monthlyDataLimit, &usedData, &requestCount, &lastRequestTime)
    if err != nil {
        logs.Logger.Error("Failed to retrieve user data", zap.Error(err), zap.Int("userID", userID))
        return err
    }

    // بررسی محدودیت حجم داده ماهانه
    if usedData+dataSize > monthlyDataLimit {
        logs.Logger.Warn("User exceeded monthly data limit", zap.Int("userID", userID), zap.Int("usedData", usedData), zap.Int("dataSize", dataSize))
        return errors.New("user has exceeded monthly data limit")
    }

    // بررسی محدودیت تعداد درخواست در دقیقه
    now := time.Now()
    if now.Sub(lastRequestTime) < time.Minute {
        if requestCount+1 > requestLimitPerMinute {
            logs.Logger.Warn("User exceeded request limit per minute", zap.Int("userID", userID), zap.Int("requestCount", requestCount))
            return errors.New("user has exceeded request limit per minute")
        }
        requestCount++
    } else {
        // ریست کردن تعداد درخواست‌ها در دقیقه اگر بیش از یک دقیقه گذشته باشد
        requestCount = 1
    }

    // بروزرسانی داده‌های کاربر
    updateQuery := `UPDATE users SET used_data = ?, request_count = ?, last_request_time = ? WHERE id = ?`
    _, err = us.db.Exec(updateQuery, usedData+dataSize, requestCount, now, userID)
    if err != nil {
        logs.Logger.Error("Failed to update user data", zap.Error(err), zap.Int("userID", userID))
        return err
    }

    logs.Logger.Info("User quota check passed", zap.Int("userID", userID), zap.Int("newUsedData", usedData+dataSize))
    return nil
}
