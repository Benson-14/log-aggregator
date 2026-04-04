package storage

import (
	"sync"
	"time"

	"github.com/Benson-14/log-aggregator/parser"
)

type MemoryStorage struct {
	mu      sync.RWMutex
	entries []*parser.LogEntry
	index   Index
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		entries: make([]*parser.LogEntry, 0),
		index:   *NewIndex(),
	}
}

func (s *MemoryStorage) Append(entry *parser.LogEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := len(s.entries)
	s.entries = append(s.entries, entry)
	s.index.Add(id, entry.Message)
	return nil

}

func (s *MemoryStorage) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

func (s *MemoryStorage) All() []*parser.LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*parser.LogEntry, len(s.entries))
	copy(result, s.entries)
	return result
}

func (s *MemoryStorage) QueryTimeRange(start, end time.Time) []*parser.LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*parser.LogEntry
	for _, entry := range s.entries {
		if entry.Timestamp.After(start) && entry.Timestamp.Before(end) {
			results = append(results, entry)
		}
	}
	return results
}

func (s *MemoryStorage) QueryLevel(level string) []*parser.LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*parser.LogEntry
	for _, entry := range s.entries {
		if entry.Level == level {
			results = append(results, entry)
		}
	}
	return results
}

func (s *MemoryStorage) Search(text string) []*parser.LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := s.index.Search(text) // index has its own internal RLock

	var results []*parser.LogEntry
	for _, id := range ids {
		if id < len(s.entries) {
			results = append(results, s.entries[id])
		}
	}
	return results
}

func (s *MemoryStorage) CountByLevel() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	counts := make(map[string]int)
	for _, entry := range s.entries {
		counts[entry.Level]++
	}
	return counts
}

func (s *MemoryStorage) CountBySource() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	counts := make(map[string]int)
	for _, entry := range s.entries {
		counts[entry.Source]++
	}
	return counts
}
