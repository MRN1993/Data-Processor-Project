package repository

import (
    "database/sql"
    "time"
    "data-processor-project/internal/domain/models"
)

type RequestRepository struct {
    db *sql.DB
}

func NewRequestRepository(db *sql.DB) *RequestRepository {
    return &RequestRepository{db: db}
}

func (repo *RequestRepository) AddRequest(request models.Request) error {
    _, err := repo.db.Exec("INSERT INTO requests (id, user_id, data, received_at) VALUES (?, ?, ?, ?)",
        request.ID, request.UserID, request.Data, request.ReceivedAt.Format(time.RFC3339))
    return err
}

func (repo *RequestRepository) GetRequestsByUserID(userID string) ([]models.Request, error) {
    rows, err := repo.db.Query("SELECT id, user_id, data, received_at FROM requests WHERE user_id = ?", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var requests []models.Request
    for rows.Next() {
        var request models.Request
        var receivedAt string
        if err := rows.Scan(&request.ID, &request.UserID, &request.Data, &receivedAt); err != nil {
            return nil, err
        }
        request.ReceivedAt, _ = time.Parse(time.RFC3339, receivedAt)
        requests = append(requests, request)
    }
    return requests, nil
}
