package models

// User represents a system user.
type User struct {
	ID     string // Unique identifier for the user
	Quota  int    // Maximum number of requests allowed
}