package logic

import (
    "errors"
    "data-processor-project/internal/domain/models"
)

// ValidateRequest validates the incoming request data.
func ValidateRequest(userID string, data string) error {
    if userID == "" || data == "" {
        return errors.New("userID and data cannot be empty")
    }
    return nil
}

// CheckDuplicate checks if the request is a duplicate.
func CheckDuplicate(requests map[string]models.Request, userID string, data string) bool {
    for _, req := range requests {
        if req.UserID == userID && req.Data == data {
            return true
        }
    }
    return false
}

// CheckUserQuota checks if the user has exceeded their request quota.
func CheckUserQuota(users map[string]models.User, requests map[string]models.Request, userID string) error {
    user, exists := users[userID]
    if !exists {
        return errors.New("user does not exist")
    }

    count := 0
    for _, req := range requests {
        if req.UserID == userID {
            count++
        }
    }

    if count >= user.Quota {
        return errors.New("user has exceeded request quota")
    }
    return nil
}