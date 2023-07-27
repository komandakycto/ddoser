package types

import "time"

// TimeSeriesIpDuplicates is a data structure that stores the count of IP addresses.
// The main aim of this data structure is to detect duplicate IP addresses within a time slot.
type TimeSeriesIpDuplicates struct {
	startTimestamp time.Time
	counter        map[string]int
	timeslot       time.Duration
	threshold      int
	duplicatesChan chan string
}

// NewTimeSeriesIpDuplicates initializes and returns a new TimeSeriesIpDuplicates data structure.
func NewTimeSeriesIpDuplicates(startTimestamp time.Time, duration time.Duration, threshold int, duplicatesChan chan string) *TimeSeriesIpDuplicates {
	return &TimeSeriesIpDuplicates{
		// initialize the start timestamp for time series logs.
		// This is used to determine the current time slot and the next time slots.
		// In general, the start timestamp is the timestamp of the first log entry.
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

// Add adds a new IP address and its timestamp to the TimeSeriesIpDuplicates data structure.
// If the timestamp is more than ts.timeslot seconds from the start timestamp, a new time slot is started.
// If the number of occurrences of an IP address in the ts.timeslot is more than ts.threshold
// the IP address is sent to the duplicates channel.
func (ts *TimeSeriesIpDuplicates) Add(ip string, timestamp time.Time) {
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
