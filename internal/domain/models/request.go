package models

import "time"

// Request represents a user request.
type Request struct {
    ID         string    `json:"id"`
    UserID     string    `json:"user_id"`
    Data       string    `json:"data"`
    ReceivedAt time.Time `json:"received_at"`
}
