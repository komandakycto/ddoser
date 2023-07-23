package main

import "time"

type TimeSeriesIpDuplicates struct {
	startTimestamp time.Time
	ipCountMap     map[string]int
	duplicatesChan chan string
	timeslot       time.Duration
	threshold      int
}

// NewTimeSeriesIpDuplicates initializes and returns a new TimeSeriesIpDuplicates data structure.
func NewTimeSeriesIpDuplicates(startTimestamp time.Time, duplicatesChan chan string, duration time.Duration, threshold int) *TimeSeriesIpDuplicates {
	return &TimeSeriesIpDuplicates{
		// initialize the start timestamp for time series logs.
		startTimestamp: startTimestamp,
		// initialize the map to store the count of IP addresses.
		ipCountMap: make(map[string]int),
		// initialize the channel to send the IP addresses that exceed the threshold.
		duplicatesChan: duplicatesChan,
		// initialize the time slot for calculation duplicates.
		timeslot: duration,
		// initialize the threshold for calculation duplicates.
		threshold: threshold,
	}
}

// Add adds a new IP address and its timestamp to the TimeSeriesIpDuplicates data structure.
// If the timestamp is more than 10 seconds from the start timestamp, a new time slot is started.
// If the number of occurrences of an IP address in the 10-second time slot is more than 10,
// the IP address is sent to the duplicates channel.
func (ts *TimeSeriesIpDuplicates) Add(ip string, timestamp time.Time) {
	if timestamp.Sub(ts.startTimestamp) > ts.timeslot {
		// Move to the next time slot of 10 seconds.
		ts.startTimestamp = timestamp
		ts.ipCountMap = make(map[string]int)
	}

	// Increment the count for the IP address within the current time slot.
	ts.ipCountMap[ip]++

	// Check if the count of the IP address exceeds the threshold (10).
	if ts.ipCountMap[ip] == ts.threshold {
		// Send the IP address to the duplicates channel.
		ts.duplicatesChan <- ip
	}
}

type LogEntry struct {
	IPAddress    string
	RequestedURL string
	Timestamp    time.Time
	UserAgent    string
}
