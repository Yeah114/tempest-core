package app

import "time"

// Message describes a general status update published by the gateway.
type Message struct {
	Type      string
	Message   string
	Error     string
	Timestamp time.Time
}
