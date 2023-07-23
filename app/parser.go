package main

import (
	"fmt"
	"regexp"
	"time"
)

const TimeLayout = "02/Jan/2006:15:04:05 -0700"

// LogEntry is a struct that represents a log entry.
func parseLogLine(logLine string) (*LogEntry, error) {
	// Regular expression to extract the required fields from the log line.
	re := regexp.MustCompile(`^([\d.]+) - - \[(\d{2}/\w+/\d{4}:\d{2}:\d{2}:\d{2} \+\d{4})\] "([A-Z]+) (.+) HTTP\/\d\.\d" \d+ \d+ "(.+)" "(.+)".*$`)

	matches := re.FindStringSubmatch(logLine)
	if len(matches) != 7 {
		return nil, fmt.Errorf("failed to parse log line: %s", logLine)
	}

	// Parse the timestamp using the given layout.
	timestamp, err := time.Parse(TimeLayout, matches[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %v", err)
	}

	entry := LogEntry{
		IPAddress:    matches[1],
		RequestedURL: matches[4],
		Timestamp:    timestamp,
		UserAgent:    matches[6],
	}

	return &entry, nil
}
