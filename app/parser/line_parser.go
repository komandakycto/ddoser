package parser

import "github.com/komandakycto/ddoser/app/entities"

type LineParser interface {
	Parse(line string) (*entities.LogEntry, error)
}
