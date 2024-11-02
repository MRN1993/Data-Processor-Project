package services

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserService struct {
    db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
    return &UserService{db: db}
}


// RegisterUser برای ثبت کاربر جدید
func (us *UserService) RegisterUser(userID int, quota int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := `INSERT INTO users (id, request_limit_per_minute, monthly_data_limit, used_data, request_count, last_request_time) VALUES (?, ?, ?, ?, ?, ?)`
    _, err := us.db.ExecContext(ctx, query, userID, quota, 0, 0, 0, time.Now())
    if err != nil {
        return errors.New("failed to register user: " + err.Error())
    }
    return nil
}


// CheckQuota بررسی محدودیت‌های کاربر
func (us *UserService) CheckQuota(userID int, dataSize int) error {
    var (
        requestLimitPerMinute, monthlyDataLimit, usedData, requestCount int
        lastRequestTime time.Time
    )

    query := `SELECT request_limit_per_minute, monthly_data_limit, used_data, request_count, last_request_time FROM users WHERE id = ?`
    row := us.db.QueryRow(query, userID)
    err := row.Scan(&requestLimitPerMinute, &monthlyDataLimit, &usedData, &requestCount, &lastRequestTime)
    if err != nil {
        return err
    }

    // بررسی محدودیت حجم داده ماهانه
    if usedData+dataSize > monthlyDataLimit {
        return errors.New("user has exceeded monthly data limit")
    }

    // بررسی محدودیت تعداد درخواست در دقیقه
    now := time.Now()
    if now.Sub(lastRequestTime) < time.Minute {
        if requestCount+1 > requestLimitPerMinute {
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
        return err
    }

    return nil
}


