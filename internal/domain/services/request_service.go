package services

import (
    "data-processor-project/internal/domain/models"
    "data-processor-project/internal/logs"
    "go.uber.org/zap"
)

type RequestRepository interface {
    AddRequest(request models.Request) error
    GetRequestsByUserID(userID int) ([]models.Request, error)
}

type RequestService struct {
    repo   RequestRepository
}

func NewRequestService(repo RequestRepository) *RequestService {
    return &RequestService{repo: repo}
}

func (s *RequestService) ProcessRequest(request models.Request) error {
    logs.Logger.Info("Processing request", zap.String("ID", request.ID))
    return s.repo.AddRequest(request)
}
