package repository

import (
    "database/sql"
    "data-processor-project/internal/domain/models"
    "data-processor-project/internal/logs"
    "go.uber.org/zap"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (repo *UserRepository) AddUser(user models.User) error {
    logs.Logger.Info("Adding new user", zap.Int("userID", user.ID), zap.Int("quota", user.Quota))

    _, err := repo.db.Exec("INSERT INTO users (id, quota) VALUES (?, ?)", user.ID, user.Quota)
    if err != nil {
        logs.Logger.Error("Failed to add user", zap.Error(err), zap.Int("userID", user.ID))
        return err
    }

    logs.Logger.Info("User added successfully", zap.Int("userID", user.ID))
    return nil
}


func (repo *UserRepository) GetUser(userID int) (models.User, error) {
    logs.Logger.Info("Fetching user data", zap.Int("userID", userID))

    var user models.User
    row := repo.db.QueryRow("SELECT id, quota FROM users WHERE id = ?", userID)
    err := row.Scan(&user.ID, &user.Quota)
    if err != nil {
        if err == sql.ErrNoRows {
            logs.Logger.Warn("User not found", zap.Int("userID", userID))
        } else {
            logs.Logger.Error("Failed to retrieve user data", zap.Error(err), zap.Int("userID", userID))
        }
        return user, err
    }

    logs.Logger.Info("User data retrieved successfully", zap.Int("userID", userID), zap.Int("quota", user.Quota))
    return user, nil
}
