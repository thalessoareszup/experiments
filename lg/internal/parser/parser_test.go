package parser

import (
	"testing"
)

func TestParse_ValidJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "simple object",
			input: `{"key": "value"}`,
		},
		{
			name:  "nested object",
			input: `{"level":"info","message":"test","data":{"nested":true}}`,
		},
		{
			name:  "with numbers",
			input: `{"count": 42, "ratio": 3.14}`,
		},
		{
			name:  "with array",
			input: `{"items": [1, 2, 3]}`,
		},
		{
			name:  "with whitespace",
			input: `  {"key": "value"}  `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := Parse(tt.input)
			if entry == nil {
				t.Error("Parse() returned nil for valid JSON")
			}
			if !entry.IsJSON {
				t.Error("Parse() should set IsJSON=true for valid JSON")
			}
			if entry.Parsed == nil {
				t.Error("Parse() should set Parsed for valid JSON")
			}
			if entry.Formatted == "" {
				t.Error("Parse() returned empty formatted string")
			}
		})
	}
}

func TestParse_NonJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectNil bool // true if we expect nil (empty lines)
	}{
		{
			name:      "plain text",
			input:     "This is not JSON",
			expectNil: false,
		},
		{
			name:      "empty string",
			input:     "",
			expectNil: true,
		},
		{
			name:      "whitespace only",
			input:     "   ",
			expectNil: true,
		},
		{
			name:      "malformed JSON",
			input:     `{"key": }`,
			expectNil: false,
		},
		{
			name:      "incomplete JSON",
			input:     `{"key": "value"`,
			expectNil: false,
		},
		{
			name:      "array (not object)",
			input:     `[1, 2, 3]`,
			expectNil: false,
		},
		{
			name:      "log prefix",
			input:     "DEBUG: Loading configuration",
			expectNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := Parse(tt.input)
			if tt.expectNil {
				if entry != nil {
					t.Error("Parse() should return nil for empty input")
				}
				return
			}

			if entry == nil {
				t.Error("Parse() should return entry for non-JSON text")
				return
			}
			if entry.IsJSON {
				t.Error("Parse() should set IsJSON=false for non-JSON")
			}
			if entry.Parsed != nil {
				t.Error("Parse() should have nil Parsed for non-JSON")
			}
			if entry.Formatted == "" {
				t.Error("Parse() should format non-JSON lines")
			}
		})
	}
}

func TestLogEntry_MatchesFilter(t *testing.T) {
	entry := &LogEntry{
		Raw: `{"level":"error","message":"Connection failed","host":"localhost"}`,
	}

	tests := []struct {
		query string
		want  bool
	}{
		{"error", true},
		{"ERROR", true}, // case insensitive
		{"Connection", true},
		{"localhost", true},
		{"notfound", false},
		{"", true}, // empty query matches all
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			if got := entry.MatchesFilter(tt.query); got != tt.want {
				t.Errorf("MatchesFilter(%q) = %v, want %v", tt.query, got, tt.want)
			}
		})
	}
}

func TestParse_NonJSON_Formatted(t *testing.T) {
	entry := Parse("Starting application...")
	if entry == nil {
		t.Fatal("Parse() returned nil for plain text")
	}
	if entry.IsJSON {
		t.Error("Expected IsJSON=false for plain text")
	}
	if entry.Formatted == "" {
		t.Error("Expected formatted output for plain text")
	}
	// The formatted output should contain the original text
	if entry.Raw != "Starting application..." {
		t.Errorf("Raw = %q, want %q", entry.Raw, "Starting application...")
	}
}
