package entities

import "time"

// LogEntry is a struct representing a parsed log entry from any log provider.
// It contains only the fields utilized in the application logic.
type LogEntry struct {
	// The data structure TimeSeriesIpDuplicates represents log rows as time series data.
	// Therefore, it is essential to always provide a value for Timestamp.
	Timestamp time.Time
	// IPAddress the main object of analysis for detecting DDoS attacks.
	IPAddress string
	// RequestedURL used to filter log rows for analysis.
	RequestedURL string
	// UserAgent used to extended analysis based on user agent data.
	UserAgent string
}
