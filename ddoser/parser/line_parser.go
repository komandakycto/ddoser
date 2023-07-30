package parser

import "github.com/komandakycto/ddoser/ddoser/entities"

// LineParser is the interface that any log provider must implement.
type LineParser interface {
	// Parse function is responsible for parsing one line of log into an inner struct.
	Parse(line string) (*entities.LogEntry, error)
}
