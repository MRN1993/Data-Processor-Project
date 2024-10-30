package events

type RequestReceived struct {
	RequestID string
}

type DuplicateDataDetected struct {
	RequestID string
}

type UserLimitExceeded struct {
	UserID string
}
