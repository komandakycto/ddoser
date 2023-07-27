package entities

import "time"

// LogEntry is a struct that represents a parsed log entry.
type LogEntry struct {
	Timestamp    time.Time
	IPAddress    string
	RequestedURL string
	UserAgent    string
}
