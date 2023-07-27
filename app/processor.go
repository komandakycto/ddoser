package main

import (
	"context"
	"github.com/komandakycto/ddoser/app/entities"
	"github.com/komandakycto/ddoser/app/parser"
	"github.com/komandakycto/ddoser/app/types"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"sync"
	"time"
)

// processGroupsConcurrently processes groups concurrently using goroutines.
func processGroupsConcurrently(ctx context.Context, groups [][]string, ipNumbersThreshold int, timeWindow int, urlPattern string, jsonLog bool, onlyIPv4 bool, log *logrus.Entry) map[string]bool {
	resultCh := make(chan string, len(groups)) // Create a channel to receive the results from each group.
	go func() {
		<-ctx.Done()
		log.Info("Context is done. Closing resultCh...")
		close(resultCh)
	}()

	var wg sync.WaitGroup
	wg.Add(len(groups))

	for i, group := range groups {
		go func(ctx context.Context, groupID int, group []string) {
			defer wg.Done()
			processGroup(ctx, groupID, group, resultCh, ipNumbersThreshold, timeWindow, urlPattern, jsonLog, log)
		}(ctx, i+1, group)
	}

	// Create a map and a mutex to store the results.
	resultMap := make(map[string]bool)
	var mu sync.Mutex

	go func() {
		for result := range resultCh {
			// Use the mutex to protect the map from concurrent writes.
			if onlyIPv4 && !isIPv4(result) {
				continue
			}

			mu.Lock()
			resultMap[result] = true
			mu.Unlock()
		}
	}()

	wg.Wait()
	close(resultCh) // Close the channel after all goroutines are done processing.

	log.Info("Finished processing all groups and writing to the file.")

	return resultMap
}

// processGroup is a function that processes each group.
func processGroup(ctx context.Context, groupID int, group []string, ch chan string, ipNumbersThreshold int, timeWindow int, urlPattern string, jsonLog bool, log *logrus.Entry) {
	log.Infof("Processing Group %d with %d elements...", groupID, len(group))

	// Calculate the initial time window start time and end time.
	firstEntry, err := parse(group[0], jsonLog)
	if err != nil {
		log.WithError(err).Error("Error parsing log line")
		return
	}

	ts := types.NewTimeSeriesIpDuplicates(
		firstEntry.Timestamp,
		time.Duration(timeWindow)*time.Second,
		ipNumbersThreshold,
		ch,
	)

	for _, element := range group {
		select {
		case <-ctx.Done():
			log.Info("Context is done. Exiting processGroup...")
			return
		default:
			// Continue processing.
		}

		// Simulate some processing time
		entry, err := parse(element, jsonLog)
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{"element": element}).Error("Error parsing log line")
			continue
		}

		// Check if the URL pattern matches or if it's an empty string (meaning all URL patterns are considered).
		if urlPattern == "" || strings.Contains(entry.RequestedURL, urlPattern) {
			ts.Add(entry.IPAddress, entry.Timestamp)
		}
	}
}

func parse(logLine string, jsonLog bool) (*entities.LogEntry, error) {
	if jsonLog {
		return parser.parseJson(logLine)
	}

	return parser.parseLogLine(logLine)
}

func isIPv4(address string) bool {
	ip := net.ParseIP(address)
	return ip != nil && ip.To4() != nil
}
