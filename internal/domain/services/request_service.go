package services

import (
	"errors"
	"fmt"
	"time"
	"data-processor-project/internal/domain/models"
	"data-processor-project/internal/generator" 
)


// RequestService provides functionality to manage requests.
type RequestService struct {
	users map[string]*models.User 
	requests map[string]models.Request 
}

// NewRequestService creates a new RequestService.
func NewRequestService() *RequestService {
	return &RequestService{
		users:    make(map[string]*models.User),
		requests: make(map[string]models.Request),
	}
}

// ValidateRequest validates the incoming request data.
func (s *RequestService) ValidateRequest(userID string, data string) error {
	if userID == "" || data == "" {
		return errors.New("userID and data cannot be empty")
	}
	return nil
}

// CheckDuplicate checks if the request is a duplicate.
func (s *RequestService) CheckDuplicate(userID string, data string) bool {
	for _, req := range s.requests {
		if req.UserID == userID && req.Data == data {
			return true
		}
	}
	return false
}

// CheckUserQuota checks if the user has exceeded their request quota.
func (s *RequestService) CheckUserQuota(userID string) bool {
	quota := s.users[userID].Quota
	count := 0
	for _, req := range s.requests {
		if req.UserID == userID {
			count++
		}
	}
	return count < quota
}

// HandleRequest processes a new request.
func (s *RequestService) HandleRequest(userID string, data string) error {
	// 1. Validate the request
	if err := s.ValidateRequest(userID, data); err != nil {
		return err
	}

	// 2. Check for duplicates
	if s.CheckDuplicate(userID, data) {
		return errors.New("duplicate request detected")
	}

	// 3. Check user quota
	if !s.CheckUserQuota(userID) {
		return errors.New("user has exceeded request quota")
	}

	// 4. Create the request
	request := models.Request{
		ID:         generator.GenerateUUID(), // Use the UUID generator
		UserID:     userID,
		Data:       data,
		ReceivedAt: time.Now(),
	}

	// Store the request in memory
	s.requests[request.ID] = request

	// Publish the appropriate event here (not shown)
	fmt.Println("Processing request:", request)
	return nil
}
