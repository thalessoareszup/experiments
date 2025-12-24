package parser

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// LogEntry represents a parsed log entry (JSON or plain text)
type LogEntry struct {
	Raw       string         // Original line
	Parsed    map[string]any // Parsed JSON data (nil for non-JSON)
	Formatted string         // Pretty-printed and colorized output
	IsJSON    bool           // Whether the entry is valid JSON
}

// Styles for JSON colorization
var (
	keyStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("81"))  // Cyan
	stringStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))  // Green
	numberStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("214")) // Orange
	boolStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // Pink
	nullStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241")) // Gray
	braceStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("245")) // Light gray
	plainTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true) // Dimmed for non-JSON
)

// Parse attempts to parse a line as JSON and returns a LogEntry
// For valid JSON, returns a pretty-printed colorized entry
// For non-JSON, returns a dimmed plain text entry
func Parse(line string) *LogEntry {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(line), &parsed); err != nil {
		// Non-JSON line - return dimmed plain text
		return &LogEntry{
			Raw:       line,
			Parsed:    nil,
			Formatted: plainTextStyle.Render(line),
			IsJSON:    false,
		}
	}

	formatted := formatJSON(parsed, 0)

	return &LogEntry{
		Raw:       line,
		Parsed:    parsed,
		Formatted: formatted,
		IsJSON:    true,
	}
}

// formatJSON recursively formats and colorizes JSON
func formatJSON(data any, indent int) string {
	indentStr := strings.Repeat("  ", indent)
	nextIndent := strings.Repeat("  ", indent+1)

	switch v := data.(type) {
	case map[string]any:
		if len(v) == 0 {
			return braceStyle.Render("{}")
		}

		var b strings.Builder
		b.WriteString(braceStyle.Render("{"))
		b.WriteString("\n")

		// Sort keys for consistent output
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i, k := range keys {
			b.WriteString(nextIndent)
			b.WriteString(keyStyle.Render("\"" + k + "\""))
			b.WriteString(": ")
			b.WriteString(formatJSON(v[k], indent+1))
			if i < len(keys)-1 {
				b.WriteString(",")
			}
			b.WriteString("\n")
		}

		b.WriteString(indentStr)
		b.WriteString(braceStyle.Render("}"))
		return b.String()

	case []any:
		if len(v) == 0 {
			return braceStyle.Render("[]")
		}

		var b strings.Builder
		b.WriteString(braceStyle.Render("["))
		b.WriteString("\n")

		for i, item := range v {
			b.WriteString(nextIndent)
			b.WriteString(formatJSON(item, indent+1))
			if i < len(v)-1 {
				b.WriteString(",")
			}
			b.WriteString("\n")
		}

		b.WriteString(indentStr)
		b.WriteString(braceStyle.Render("]"))
		return b.String()

	case string:
		return stringStyle.Render("\"" + escapeString(v) + "\"")

	case float64:
		return numberStyle.Render(formatNumber(v))

	case bool:
		if v {
			return boolStyle.Render("true")
		}
		return boolStyle.Render("false")

	case nil:
		return nullStyle.Render("null")

	default:
		return stringStyle.Render(formatAny(v))
	}
}

// escapeString escapes special characters in a string
func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// formatNumber formats a number, removing trailing zeros for integers
func formatNumber(n float64) string {
	if n == float64(int64(n)) {
		return strings.TrimSuffix(strings.TrimSuffix(
			strings.ReplaceAll(formatAny(int64(n)), ",", ""), ".0"), ".00")
	}
	return formatAny(n)
}

// formatAny formats any value as a string
func formatAny(v any) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
	return strings.TrimSpace(buf.String())
}

// MatchesFilter checks if the log entry matches a search query
func (e *LogEntry) MatchesFilter(query string) bool {
	if query == "" {
		return true
	}
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(e.Raw), query)
}
