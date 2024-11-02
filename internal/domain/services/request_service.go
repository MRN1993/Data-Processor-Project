package services

import (
    "data-processor-project/internal/domain/models"
    "go.uber.org/zap"
)

type RequestRepository interface {
    AddRequest(request models.Request) error
    GetRequestsByUserID(userID int) ([]models.Request, error)
}

type RequestService struct {
    repo   RequestRepository
    logger *zap.Logger
}

func NewRequestService(repo RequestRepository, logger *zap.Logger) *RequestService {
    return &RequestService{repo: repo, logger: logger}
}

func (s *RequestService) ProcessRequest(request models.Request) error {
    s.logger.Info("Processing request", zap.String("ID", request.ID))
    return s.repo.AddRequest(request)
}
