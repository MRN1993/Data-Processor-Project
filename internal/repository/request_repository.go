package repository

import (
    "database/sql"
    "data-processor-project/internal/domain/models"
    "time"
)

type SQLRequestRepository struct {
    db *sql.DB
}

func NewSQLRequestRepository(db *sql.DB) *SQLRequestRepository {
    return &SQLRequestRepository{db: db}
}

func (r *SQLRequestRepository) AddRequest(request models.Request) error {
    query := `INSERT INTO requests (id, user_id, data, received_at) VALUES (?, ?, ?, ?)`
    _, err := r.db.Exec(query, request.ID, request.UserID, request.Data, time.Now())
    return err
}

func (r *SQLRequestRepository) GetRequestsByUserID(userID int) ([]models.Request, error) {
    query := `SELECT id, user_id, data, received_at FROM requests WHERE user_id = ?`
    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var requests []models.Request
    for rows.Next() {
        var request models.Request
        if err := rows.Scan(&request.ID, &request.UserID, &request.Data, &request.ReceivedAt); err != nil {
            return nil, err
        }
        requests = append(requests, request)
    }
    return requests, nil
}
