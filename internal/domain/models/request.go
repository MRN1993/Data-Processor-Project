package models

import "time"

type Request struct {
    ID         string
    UserID     string
    Data       string
    ReceivedAt time.Time
}
