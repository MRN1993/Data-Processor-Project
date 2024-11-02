package services

import (
    "data-processor-project/internal/domain/logic"
    "data-processor-project/internal/logs"
    "database/sql"

    "go.uber.org/zap"
)

type UserService struct {
    db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
    return &UserService{db: db}
}

// RegisterUser registers a new user in the database.
func (us *UserService) RegisterUser(quota, monthlyDataLimit, requestLimitPerMinute int) error {
    logs.Logger.Info("Registering new user", zap.Int("quota", quota))

    err := logic.RegisterUserInDB(us.db, quota, monthlyDataLimit, requestLimitPerMinute)
    if err != nil {
        logs.Logger.Error("Failed to register user", zap.Error(err))
        return err
    }

    logs.Logger.Info("User registered successfully")
    return nil
}

