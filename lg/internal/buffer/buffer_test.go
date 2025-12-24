package buffer

import (
	"testing"

	"github.com/thalessoares/lg/internal/parser"
)

func TestBuffer_AddAndGet(t *testing.T) {
	buf := New(10)

	entry := &parser.LogEntry{Raw: `{"test": 1}`}
	buf.Add(entry)

	if buf.Len() != 1 {
		t.Errorf("Len() = %d, want 1", buf.Len())
	}

	got := buf.Get(0)
	if got != entry {
		t.Error("Get(0) did not return the added entry")
	}
}

func TestBuffer_Capacity(t *testing.T) {
	capacity := 5
	buf := New(capacity)

	// Add more entries than capacity
	for i := 0; i < 10; i++ {
		buf.Add(&parser.LogEntry{Raw: `{"i": ` + string(rune('0'+i)) + `}`})
	}

	if buf.Len() != capacity {
		t.Errorf("Len() = %d, want %d (capacity)", buf.Len(), capacity)
	}
}

func TestBuffer_Entries(t *testing.T) {
	buf := New(10)

	entries := []*parser.LogEntry{
		{Raw: `{"a": 1}`},
		{Raw: `{"b": 2}`},
		{Raw: `{"c": 3}`},
	}

	for _, e := range entries {
		buf.Add(e)
	}

	got := buf.Entries()
	if len(got) != len(entries) {
		t.Errorf("Entries() len = %d, want %d", len(got), len(entries))
	}

	// Verify it's a copy
	got[0] = nil
	if buf.Get(0) == nil {
		t.Error("Entries() should return a copy, not the original slice")
	}
}

func TestBuffer_Filter(t *testing.T) {
	buf := New(10)

	buf.Add(&parser.LogEntry{Raw: `{"level": "error", "message": "failed"}`})
	buf.Add(&parser.LogEntry{Raw: `{"level": "info", "message": "success"}`})
	buf.Add(&parser.LogEntry{Raw: `{"level": "error", "message": "timeout"}`})

	filtered := buf.Filter("error")
	if len(filtered) != 2 {
		t.Errorf("Filter('error') len = %d, want 2", len(filtered))
	}

	filtered = buf.Filter("success")
	if len(filtered) != 1 {
		t.Errorf("Filter('success') len = %d, want 1", len(filtered))
	}

	filtered = buf.Filter("")
	if len(filtered) != 3 {
		t.Errorf("Filter('') len = %d, want 3", len(filtered))
	}
}

func TestBuffer_Clear(t *testing.T) {
	buf := New(10)

	buf.Add(&parser.LogEntry{Raw: `{"a": 1}`})
	buf.Add(&parser.LogEntry{Raw: `{"b": 2}`})

	buf.Clear()

	if buf.Len() != 0 {
		t.Errorf("Len() after Clear() = %d, want 0", buf.Len())
	}
}

func TestBuffer_GetOutOfBounds(t *testing.T) {
	buf := New(10)
	buf.Add(&parser.LogEntry{Raw: `{"a": 1}`})

	if got := buf.Get(-1); got != nil {
		t.Error("Get(-1) should return nil")
	}

	if got := buf.Get(100); got != nil {
		t.Error("Get(100) should return nil for index out of bounds")
	}
}
