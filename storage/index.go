package storage

import (
	"sort"
	"strings"
	"sync"
)

type Index struct {
	mu    sync.RWMutex
	words map[string][]int
}

func NewIndex() *Index {
	return &Index{
		words: make(map[string][]int),
	}
}

func tokenize(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func (idx *Index) Add(entryID int, text string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	words := tokenize(text)

	for _, word := range words {
		idx.words[word] = append(idx.words[word], entryID)
	}
}

// Search returns entry IDs containing the word (case-insensitive).
// A copy of the internal slice is returned so callers cannot accidentally
// mutate the index's backing array.
func (idx *Index) Search(word string) []int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	word = strings.ToLower(word)
	entryIDs, ok := idx.words[word]
	if !ok {
		return nil
	}
	// Return a copy — not a direct alias into the internal map.
	result := make([]int, len(entryIDs))
	copy(result, entryIDs)
	return result
}

// SearchAll returns entry IDs containing ALL words (AND search)
func (idx *Index) SearchAll(words []string) []int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if len(words) == 0 {
		return nil
	}

	result := make(map[int]bool)
	firstWord := strings.ToLower(words[0])
	for _, id := range idx.words[firstWord] {
		result[id] = true
	}

	for _, word := range words[1:] {
		word = strings.ToLower(word)
		newResult := make(map[int]bool)
		for _, id := range idx.words[word] {
			if result[id] {
				newResult[id] = true
			}
		}
		result = newResult
	}

	var ids []int
	for id := range result {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	return ids
}
