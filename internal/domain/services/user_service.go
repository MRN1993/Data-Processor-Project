package services

import (
    "errors"
    "data-processor-project/internal/domain/models"
    "data-processor-project/internal/generator"
)

// UserService manages users in the system.
type UserService struct {
    users map[string]models.User // Store users
}

// NewUserService creates a new UserService.
func NewUserService() *UserService {
    return &UserService{
        users: make(map[string]models.User),
    }
}

// CreateUser creates a new user with a specified quota.
func (s *UserService) CreateUser(quota int) (models.User, error) {
    if quota <= 0 {
        return models.User{}, errors.New("quota must be greater than zero")
    }

    userID := generator.GenerateUUID()
    user := models.User{
        ID:    userID,
        Quota: quota,
    }

    // Store the user
    s.users[userID] = user
    return user, nil
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(userID string) (models.User, error) {
    user, exists := s.users[userID]
    if !exists {
        return models.User{}, errors.New("user does not exist")
    }
    return user, nil
}
