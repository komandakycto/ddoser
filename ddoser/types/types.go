package types

import "time"

// TimeSeriesIpAnalyzer is a data structure designed to store the count of IP addresses.
// Its primary purpose is to detect duplicate IP addresses within a given time slot.
type TimeSeriesIpAnalyzer struct {
	startTimestamp time.Time
	counter        map[string]int
	timeslot       time.Duration
	threshold      int
	duplicatesChan chan string
}

// NewTimeSeriesIpAnalyzer initializes and returns a new TimeSeriesIpAnalyzer data structure.
func NewTimeSeriesIpAnalyzer(startTimestamp time.Time, duration time.Duration, threshold int, duplicatesChan chan string) *TimeSeriesIpAnalyzer {
	return &TimeSeriesIpAnalyzer{
		// Initialize the start timestamp for time series logs.
		// This timestamp is utilized to determine the current time slot and the subsequent time slots.
		// Typically, the start timestamp corresponds to the timestamp of the first log entry.
		startTimestamp: startTimestamp,
		// initialize the map to store the count of IP addresses.
		counter: make(map[string]int),
		// initialize the channel to send the IP addresses that exceed the threshold.
		duplicatesChan: duplicatesChan,
		// initialize the time slot for calculation duplicates.
		// During the time slot, the IP addresses are counted.
		// For high traffic, the time slot can be reduced to 1 second.
		// For low traffic, the time slot can be increased to 1 minute.
		timeslot: duration,
		// initialize the threshold for calculation duplicates.
		// If the number of occurrences of an IP address in the time slot is more than the threshold,
		// the IP address is sent to the duplicates channel.
		threshold: threshold,
	}
}

// Add adds a new IP address and its timestamp to the TimeSeriesIpAnalyzer data structure.
// If the timestamp is more than ts.timeslot seconds from the start timestamp, a new time slot is started.
// If the number of occurrences of an IP address in the ts.timeslot is more than ts.threshold
// the IP address is sent to the duplicates channel.
func (ts *TimeSeriesIpAnalyzer) Add(ip string, timestamp time.Time) {
	if timestamp.Sub(ts.startTimestamp) > ts.timeslot {
		// Move to the next time slot.
		ts.startTimestamp = timestamp
		ts.counter = make(map[string]int)
	}

	// Increment the count for the IP address within the current time slot.
	ts.counter[ip]++

	// Check if the count of the IP address exceeds the threshold.
	// Send ip only once to the duplicates channel.
	if ts.counter[ip] == ts.threshold {
		// Send the IP address to the duplicates channel.
		ts.duplicatesChan <- ip
	}
}
