package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thalessoares/lg/internal/buffer"
	"github.com/thalessoares/lg/internal/parser"
)

// Mode represents the current UI mode
type Mode int

const (
	ModeView Mode = iota
	ModeSearch
)

// LogMsg is sent when a new log entry is received
type LogMsg *parser.LogEntry

// Model is the main TUI model
type Model struct {
	buffer       *buffer.Buffer
	viewport     viewport.Model
	searchInput  textinput.Model
	mode         Mode
	paused       bool
	filter       string
	width        int
	height       int
	ready        bool
	autoScroll   bool
	entries      []*parser.LogEntry // Filtered entries for display
	totalEntries int                // Total entries in buffer
}

// New creates a new Model
func New(buf *buffer.Buffer) Model {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 256
	ti.Width = 50

	return Model{
		buffer:     buf,
		searchInput: ti,
		mode:       ModeView,
		paused:     false,
		autoScroll: true,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 1 // Status bar
		footerHeight := 2 // Help + search bar (when visible)

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-headerHeight-footerHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - headerHeight - footerHeight
		}
		m.updateViewportContent()

	case LogMsg:
		if msg != nil {
			m.buffer.Add(msg)
			if !m.paused {
				m.updateViewportContent()
				if m.autoScroll {
					m.viewport.GotoBottom()
				}
			}
		}
	}

	// Update viewport
	if m.mode == ModeView {
		var vpCmd tea.Cmd
		m.viewport, vpCmd = m.viewport.Update(msg)
		cmds = append(cmds, vpCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case ModeSearch:
		return m.handleSearchMode(msg)
	default:
		return m.handleViewMode(msg)
	}
}

func (m Model) handleViewMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "/":
		m.mode = ModeSearch
		m.searchInput.Focus()
		return m, textinput.Blink

	case "p":
		m.paused = !m.paused
		if !m.paused {
			m.updateViewportContent()
			if m.autoScroll {
				m.viewport.GotoBottom()
			}
		}

	case "j", "down":
		m.viewport.LineDown(1)
		m.autoScroll = m.viewport.AtBottom()

	case "k", "up":
		m.viewport.LineUp(1)
		m.autoScroll = false

	case "g":
		m.viewport.GotoTop()
		m.autoScroll = false

	case "G":
		m.viewport.GotoBottom()
		m.autoScroll = true

	case "ctrl+d", "pgdown":
		m.viewport.HalfViewDown()
		m.autoScroll = m.viewport.AtBottom()

	case "ctrl+u", "pgup":
		m.viewport.HalfViewUp()
		m.autoScroll = false

	case "c":
		m.buffer.Clear()
		m.filter = ""
		m.updateViewportContent()

	case "esc":
		if m.filter != "" {
			m.filter = ""
			m.searchInput.SetValue("")
			m.updateViewportContent()
		}
	}

	return m, nil
}

func (m Model) handleSearchMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "enter":
		m.filter = m.searchInput.Value()
		m.mode = ModeView
		m.searchInput.Blur()
		m.updateViewportContent()
		return m, nil

	case "esc":
		m.mode = ModeView
		m.searchInput.Blur()
		m.searchInput.SetValue(m.filter)
		return m, nil
	}

	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m *Model) updateViewportContent() {
	if !m.ready {
		return
	}

	m.entries = m.buffer.Filter(m.filter)
	m.totalEntries = m.buffer.Len()

	var content strings.Builder
	separator := separatorStyle.Render(strings.Repeat("─", m.width-2))

	for i, entry := range m.entries {
		content.WriteString(entry.Formatted)
		if i < len(m.entries)-1 {
			content.WriteString("\n")
			content.WriteString(separator)
			content.WriteString("\n")
		}
	}

	m.viewport.SetContent(content.String())
}

// View implements tea.Model
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var b strings.Builder

	// Status bar
	b.WriteString(m.renderStatusBar())
	b.WriteString("\n")

	// Main viewport
	b.WriteString(m.viewport.View())
	b.WriteString("\n")

	// Search bar or help
	if m.mode == ModeSearch {
		b.WriteString(m.renderSearchBar())
	} else {
		b.WriteString(m.renderHelp())
	}

	return b.String()
}

func (m Model) renderStatusBar() string {
	// Mode indicator
	var modeStr string
	if m.paused {
		modeStr = statusPausedStyle.Render("PAUSED")
	} else if m.filter != "" {
		modeStr = statusSearchStyle.Render("FILTER")
	} else {
		modeStr = statusModeStyle.Render("VIEW")
	}

	// Entry count
	countStr := statusInfoStyle.Render(
		fmt.Sprintf("Entries: %d/%d", len(m.entries), m.totalEntries),
	)

	// Filter info
	var filterStr string
	if m.filter != "" {
		filterStr = statusInfoStyle.Render(fmt.Sprintf("Filter: %q", m.filter))
	}

	// Scroll position
	scrollStr := statusInfoStyle.Render(
		fmt.Sprintf("%.0f%%", m.viewport.ScrollPercent()*100),
	)

	// Build status bar
	left := lipgloss.JoinHorizontal(lipgloss.Left, modeStr, countStr, filterStr)
	right := scrollStr

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}

	return statusBarStyle.Render(
		lipgloss.JoinHorizontal(lipgloss.Left,
			left,
			strings.Repeat(" ", gap),
			right,
		),
	)
}

func (m Model) renderSearchBar() string {
	prompt := searchPromptStyle.Render("/")
	return searchBarStyle.Render(prompt + m.searchInput.View())
}

func (m Model) renderHelp() string {
	helpItems := []string{
		"j/k: scroll",
		"g/G: top/bottom",
		"/: search",
		"p: pause",
		"c: clear",
		"q: quit",
	}
	return helpStyle.Render(strings.Join(helpItems, " | "))
}

// AddLogEntry is called to add a new log entry (for use with Program.Send)
func AddLogEntry(entry *parser.LogEntry) tea.Msg {
	return LogMsg(entry)
}
