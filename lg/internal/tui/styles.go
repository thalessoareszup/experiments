package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("81")  // Cyan
	secondaryColor = lipgloss.Color("214") // Orange
	successColor   = lipgloss.Color("82")  // Green
	warningColor   = lipgloss.Color("214") // Orange
	errorColor     = lipgloss.Color("196") // Red
	mutedColor     = lipgloss.Color("241") // Gray
	highlightColor = lipgloss.Color("226") // Yellow

	// Status bar styles
	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("252")).
			Padding(0, 1)

	statusModeStyle = lipgloss.NewStyle().
			Background(primaryColor).
			Foreground(lipgloss.Color("0")).
			Padding(0, 1).
			Bold(true)

	statusPausedStyle = lipgloss.NewStyle().
				Background(warningColor).
				Foreground(lipgloss.Color("0")).
				Padding(0, 1).
				Bold(true)

	statusSearchStyle = lipgloss.NewStyle().
				Background(successColor).
				Foreground(lipgloss.Color("0")).
				Padding(0, 1).
				Bold(true)

	statusInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Padding(0, 1)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(0, 1)

	// Search bar styles
	searchBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("252")).
			Padding(0, 1)

	searchPromptStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	// Log entry styles
	entryStyle = lipgloss.NewStyle().
			Padding(0, 1).
			MarginBottom(1)

	selectedEntryStyle = lipgloss.NewStyle().
				Padding(0, 1).
				MarginBottom(1).
				Background(lipgloss.Color("237"))

	// Separator style
	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238"))

	// Highlight style for search matches
	highlightMatchStyle = lipgloss.NewStyle().
				Background(highlightColor).
				Foreground(lipgloss.Color("0"))

	// Title style
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	// Border style
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238"))
)
