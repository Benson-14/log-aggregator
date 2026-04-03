package storage

import (
	"time"

	"github.com/Benson-14/log-aggregator/parser"
)

type Storage interface {
	// Append adds a parsed log entry to the store.
	Append(entry *parser.LogEntry) error

	// Len returns the total number of stored entries.
	Len() int

	// QueryTimeRange returns entries whose timestamp falls in [start, end).
	QueryTimeRange(start, end time.Time) []*parser.LogEntry

	// QueryLevel returns all entries whose Level matches the given string.
	QueryLevel(level string) []*parser.LogEntry

	// Search performs a full-text search over log messages and returns matching entries.
	Search(text string) []*parser.LogEntry

	// CountByLevel returns a map from log-level string to entry count.
	CountByLevel() map[string]int

	// CountBySource returns a map from source string to entry count.
	CountBySource() map[string]int
}
