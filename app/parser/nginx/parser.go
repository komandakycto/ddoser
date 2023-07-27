package nginx

import (
	"encoding/json"
	"fmt"
	"github.com/komandakycto/ddoser/app/entities"
	"regexp"
	"time"
)

// DefaultTimeLayout is the default layout of the timestamp in the nginx log.
const DefaultTimeLayout = "02/Jan/2006:15:04:05 -0700"

// logRow is an inner struct that represents a log row in json format.
// Use default json field names.
type logRow struct {
	IPAddress    string `json:"ip"`
	RequestedURL string `json:"uri"`
	Timestamp    string `json:"time"`
	UserAgent    string `json:"user_agent"`
}

// Parser is a struct that represents a nginx parser.
type Parser struct {
	// isDefaultFormat nginx log in default format or json.
	isDefaultFormat bool
	// defaultRegex precompiled regex for parsing default nginx log.
	defaultRegex *regexp.Regexp
	// jsonTimeLayout layout for time field in json nginx log.
	jsonTimeLayout string
}

// NewNginxParser is a function that creates instance of nginx parser.
func NewNginxParser(isDefaultFormat bool, jsonTimeLayout *string) *Parser {
	// Default time layout is RFC3339.
	jsonTime := time.RFC3339
	if jsonTimeLayout != nil {
		// Use the given layout.
		jsonTime = *jsonTimeLayout
	}

	var re *regexp.Regexp
	if isDefaultFormat {
		re = regexp.MustCompile(`^([\d.]+) - - \[(\d{2}/\w+/\d{4}:\d{2}:\d{2}:\d{2} \+\d{4})\] "([A-Z]+) (.+) HTTP\/\d\.\d" \d+ \d+ "(.+)" "(.+)".*$`)
	}

	return &Parser{
		isDefaultFormat: isDefaultFormat,
		defaultRegex:    re,
		jsonTimeLayout:  jsonTime,
	}
}

// Parse is a function that parses a log line.
func (p *Parser) Parse(line string) (*entities.LogEntry, error) {
	if p.isDefaultFormat {
		return p.parseDefault(line)
	}

	return p.parseJson(line)
}

func (p *Parser) parseDefault(logLine string) (*entities.LogEntry, error) {
	matches := p.defaultRegex.FindStringSubmatch(logLine)
	if len(matches) != 7 {
		return nil, fmt.Errorf("failed to parse nginix log line: %s", logLine)
	}

	// Parse the timestamp using the given layout.
	timestamp, err := time.Parse(DefaultTimeLayout, matches[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %v", err)
	}

	entry := entities.LogEntry{
		IPAddress:    matches[1],
		RequestedURL: matches[4],
		Timestamp:    timestamp,
		UserAgent:    matches[6],
	}

	return &entry, nil
}

func (p *Parser) parseJson(logLine string) (*entities.LogEntry, error) {
	var logEntry logRow
	err := json.Unmarshal([]byte(logLine), &logEntry)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}

	timestamp, err := time.Parse(p.jsonTimeLayout, logEntry.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %v", err)
	}

	return &entities.LogEntry{
		IPAddress:    logEntry.IPAddress,
		RequestedURL: logEntry.RequestedURL,
		Timestamp:    timestamp,
		UserAgent:    logEntry.UserAgent,
	}, nil
}
