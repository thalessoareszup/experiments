package buffer

import (
	"sync"

	"github.com/thalessoares/lg/internal/parser"
)

// Buffer is a thread-safe ring buffer for log entries
type Buffer struct {
	entries  []*parser.LogEntry
	capacity int
	mu       sync.RWMutex
}

// New creates a new Buffer with the specified capacity
func New(capacity int) *Buffer {
	if capacity <= 0 {
		capacity = 10000
	}
	return &Buffer{
		entries:  make([]*parser.LogEntry, 0, capacity),
		capacity: capacity,
	}
}

// Add adds a new entry to the buffer
func (b *Buffer) Add(entry *parser.LogEntry) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.entries) >= b.capacity {
		// Remove oldest entry
		b.entries = b.entries[1:]
	}
	b.entries = append(b.entries, entry)
}

// Entries returns a copy of all entries
func (b *Buffer) Entries() []*parser.LogEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]*parser.LogEntry, len(b.entries))
	copy(result, b.entries)
	return result
}

// Len returns the number of entries in the buffer
func (b *Buffer) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.entries)
}

// Get returns the entry at the specified index
func (b *Buffer) Get(index int) *parser.LogEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if index < 0 || index >= len(b.entries) {
		return nil
	}
	return b.entries[index]
}

// Filter returns entries that match the query
func (b *Buffer) Filter(query string) []*parser.LogEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if query == "" {
		result := make([]*parser.LogEntry, len(b.entries))
		copy(result, b.entries)
		return result
	}

	var result []*parser.LogEntry
	for _, entry := range b.entries {
		if entry.MatchesFilter(query) {
			result = append(result, entry)
		}
	}
	return result
}

// Clear removes all entries from the buffer
func (b *Buffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = b.entries[:0]
}
