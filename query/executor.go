package query

import (
	"strings"

	"github.com/Benson-14/log-aggregator/parser"
	"github.com/Benson-14/log-aggregator/storage"
)

// Executor runs queries against a Storage back-end.
type Executor struct {
	store storage.Storage
}

// NewExecutor returns an Executor backed by the given Storage.
func NewExecutor(store storage.Storage) *Executor {
	return &Executor{store: store}
}

// Run executes a raw query string and returns the matching log entries.
// It combines field filters (level:error, source:db) with optional
// full-text matching on the message field.
func (e *Executor) Run(input string) []*parser.LogEntry {
	q := Parse(input)

	// Start with level or source filter if present (fast path via Storage).
	var candidates []*parser.LogEntry

	switch {
	case q.Fields["level"] != "" && q.Fields["source"] != "":
		// Both filters: query by level then filter by source in-process.
		for _, entry := range e.store.QueryLevel(q.Fields["level"]) {
			if entry.Source == q.Fields["source"] {
				candidates = append(candidates, entry)
			}
		}
	case q.Fields["level"] != "":
		candidates = e.store.QueryLevel(q.Fields["level"])
	case q.Fields["source"] != "":
		// No dedicated QuerySource — iterate via Search("")+filter.
		for _, entry := range e.store.Search("") {
			if entry.Source == q.Fields["source"] {
				candidates = append(candidates, entry)
			}
		}
	default:
		// No field filters — start with all entries via full-text search.
		if q.Text != "" {
			candidates = e.store.Search(q.Text)
		}
	}

	// Apply full-text filter on candidates if text term is also present.
	if q.Text == "" {
		return candidates
	}

	lower := strings.ToLower(q.Text)
	var results []*parser.LogEntry
	for _, entry := range candidates {
		if strings.Contains(strings.ToLower(entry.Message), lower) {
			results = append(results, entry)
		}
	}
	return results
}
