package event_handlers

import (
	"data-processor-project/internal/domain/events"
	"fmt"
)

// RequestReceivedHandler processes RequestReceived events.
type RequestReceivedHandler struct{}

// Handle processes the event.
func (h *RequestReceivedHandler) Handle(event events.RequestReceived) {
	fmt.Printf("Handling Request Received Event: %s\n", event.RequestID)
	// Add logic to handle the request
}
