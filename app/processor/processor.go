package processor

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/komandakycto/ddoser/app/helpers"
	"github.com/komandakycto/ddoser/app/parser"
	"github.com/komandakycto/ddoser/app/types"
)

// IPAnalysis is the struct that analyzes the log entries and finds the IPs that exceed the threshold.
type IPAnalysis struct {
	// threshold is the number of requests that an IP must exceed in the given time window.
	threshold int
	// timeWindow is the time window in seconds.
	timeWindow int
	// url is the URL pattern to filter the log entries.
	url string
	// parser is the log parser.
	parser parser.LineParser
	// onlyIPv4 is a flag to indicate if only IPv4 addresses should be considered.
	onlyIPv4 bool
	// logger is the logger instance.
	logger *logrus.Entry
}

// NewIPAnalysis creates a new IPAnalysis instance.
func NewIPAnalysis(threshold int, timeWindow int, url string, parser parser.LineParser, onlyIPv4 bool, logger *logrus.Entry) *IPAnalysis {
	return &IPAnalysis{
		threshold:  threshold,
		timeWindow: timeWindow,
		url:        url,
		parser:     parser,
		onlyIPv4:   onlyIPv4,
		logger:     logger,
	}
}

// Process function is responsible for processing the log entries and finding the IPs that exceed the threshold.
// It returns a map of IPs that exceed the threshold.
func (a *IPAnalysis) Process(ctx context.Context, groups [][]string) (map[string]bool, error) {
	resultCh := make(chan string, len(groups)) // Create a channel to receive the results from each group.

	var (
		// Create a wait group to wait for all the goroutines to finish.
		wg sync.WaitGroup
		// Create a wait group to wait for the collector goroutine to finish.
		collectorWg sync.WaitGroup
	)

	wg.Add(len(groups))
	for i, group := range groups {
		go func(ctx context.Context, groupID int, group []string) {
			defer wg.Done()
			a.processGroup(ctx, groupID, group, resultCh)
		}(ctx, i+1, group)
	}

	// Create a map and a mutex to store the results.
	resultMap := make(map[string]bool)
	var mu sync.Mutex

	// Create a goroutine to receive the results from the channel and store them in the map.
	collectorWg.Add(1)
	go func() {
		defer collectorWg.Done()

		for result := range resultCh {
			// Use the mutex to protect the map from concurrent writes.
			if a.onlyIPv4 && !helpers.IsIPv4(result) {
				continue
			}

			mu.Lock()
			resultMap[result] = true
			mu.Unlock()
		}
	}()

	// Wait for all the goroutines to finish.
	wg.Wait()
	close(resultCh)

	// Wait for the collector goroutine to finish.
	collectorWg.Wait()

	a.logger.Info("Finished processing all groups.")

	return resultMap, nil
}

func (a *IPAnalysis) processGroup(ctx context.Context, groupID int, group []string, ch chan string) {
	a.logger.Infof("Processing Group %d with %d elements...", groupID, len(group))

	// Calculate the initial time window start time and end time.
	firstEntry, err := a.parser.Parse(group[0])
	if err != nil {
		a.logger.WithError(err).Error("Error parsing log line")
		return
	}

	ts := types.NewTimeSeriesIpAnalyzer(
		firstEntry.Timestamp,
		time.Duration(a.timeWindow)*time.Second,
		a.threshold,
		ch,
	)

	for _, element := range group {
		select {
		case <-ctx.Done():
			a.logger.Info("Context is done. Exiting processGroup...")
			return
		default:
			// Continue processing.
		}

		entry, err := a.parser.Parse(element)
		if err != nil {
			a.logger.WithError(err).WithFields(logrus.Fields{"element": element}).Error("Error parsing log line")
			continue
		}

		// Check if the URL pattern matches or if it's an empty string (meaning all URL patterns are considered).
		if a.url == "" || strings.Contains(entry.RequestedURL, a.url) {
			ts.Add(entry.IPAddress, entry.Timestamp)
		}
	}
}
