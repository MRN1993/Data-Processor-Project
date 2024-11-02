package services

import (
    "errors"
    "database/sql"
    "data-processor-project/internal/logs"
    "go.uber.org/zap"
    "data-processor-project/internal/domain/logic"
)

type UserService struct {
    db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
    return &UserService{db: db}
}

// متد RegisterUser برای ثبت یک کاربر جدید
func (us *UserService) RegisterUser(quota int, monthlyDataLimit int, requestLimitPerMinute int) error {
    logs.Logger.Info("Attempting to register user", zap.Int("quota", quota))

    if err := logic.RegisterUserInDB(us.db, quota, monthlyDataLimit, requestLimitPerMinute); err != nil {
        logs.Logger.Error("Failed to register user", zap.Error(err))
        return errors.New("failed to register user: " + err.Error())
    }

    logs.Logger.Info("User registered successfully")
    return nil
}

// متد CheckQuota برای بررسی محدودیت‌های کاربر
func (us *UserService) CheckQuota(userID int, dataSize int) error {
    logs.Logger.Info("Checking user quota", zap.Int("userID", userID), zap.Int("dataSize", dataSize))

    if err := logic.CheckUserLimits(us.db, userID, dataSize); err != nil {
        logs.Logger.Warn("Quota check failed", zap.Error(err))
        return err
    }

    logs.Logger.Info("User quota check passed", zap.Int("userID", userID))
    return nil
}
