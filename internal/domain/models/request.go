package models

import "time"

// Request represents a user request.
type Request struct {
	ID         string    // Unique identifier
	UserID     string    // ID of the user making the request
	Data       string    // The request data
	ReceivedAt time.Time // Time the request was received
}