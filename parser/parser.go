package parser

import "time"

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Source    string
	Message   string
	Fields    map[string]string
}

type Parser interface {
	Parse(line string) (*LogEntry, error)
}
