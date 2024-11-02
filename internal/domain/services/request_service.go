package services

import (
	"data-processor-project/internal/domain/logic"
	"data-processor-project/internal/logs"
	"errors"
	"database/sql"

	"go.uber.org/zap"
)

type RequestService struct {
	db *sql.DB
}

func NewRequestService(db *sql.DB) *RequestService {
	return &RequestService{db: db}
}

func (s *RequestService) ProcessRequest(UserID int, Data string) error {
	logs.Logger.Info("Processing request")

	// Validate the request data
	if err := logic.ValidateRequest(s.db, UserID, Data); err != nil {
		logs.Logger.Error("Validation failed", zap.Error(err))
		return err
	}

	// Check for duplicates
	isDuplicate, err := logic.CheckDuplicate(s.db, UserID, Data)
	if err != nil {
		logs.Logger.Error("Failed to check for duplicate request", zap.Error(err))
		return err
	}
	if isDuplicate {
		logs.Logger.Warn("Duplicate request detected")
		return errors.New("duplicate request")
	}

	// Check user quota
	if err := logic.CheckUserQuota(s.db, UserID); err != nil {
		logs.Logger.Warn("User quota check failed", zap.Error(err))
		return err
	}


    //  Check user dataSize  
    dataSize := len(Data)
    if err := logic.CheckUserLimits(s.db, UserID, dataSize); err != nil {
        logs.Logger.Warn("User limit check failed", zap.Int("userID", UserID), zap.Error(err))
        return err
    }

    userLimits, err := logic.RetrieveUserLimits(s.db, UserID)
    if err != nil {
        logs.Logger.Error("Failed to retrieve user limits", zap.Error(err))
        return err
    }

    err = logic.UpdateUserQuota(s.db, UserID, dataSize, userLimits)
    if err != nil {
        logs.Logger.Error("Failed to update user quota", zap.Error(err))
        return err
    }
	if err := logic.RegisterRequest(s.db, UserID, Data); err != nil {
		logs.Logger.Error("Failed to add request", zap.Error(err))
		return err
	}

	logs.Logger.Info("Request processed successfully")
	return nil
}
