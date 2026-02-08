// Package tui provides shared TUI components for Lango CLI commands.
package tui

import "github.com/charmbracelet/lipgloss"

// Color palette for consistent theming
var (
	Primary    = lipgloss.Color("#7C3AED") // Purple
	Success    = lipgloss.Color("#10B981") // Green
	Warning    = lipgloss.Color("#F59E0B") // Amber
	Error      = lipgloss.Color("#EF4444") // Red
	Muted      = lipgloss.Color("#6B7280") // Gray
	Foreground = lipgloss.Color("#F9FAFB") // White
	Background = lipgloss.Color("#1F2937") // Dark gray
	Highlight  = lipgloss.Color("#3B82F6") // Blue
)

// Base styles for TUI components
var (
	// TitleStyle for main headers
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	// SubtitleStyle for secondary headers
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Muted).
			MarginBottom(1)

	// SuccessStyle for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success)

	// WarningStyle for warning messages
	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	// ErrorStyle for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error)

	// MutedStyle for less important text
	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	// HighlightStyle for emphasized text
	HighlightStyle = lipgloss.NewStyle().
			Foreground(Highlight).
			Bold(true)

	// BoxStyle for bordered containers
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Muted).
			Padding(1, 2)

	// ListItemStyle for list items
	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	// SelectedItemStyle for selected list items
	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(Primary).
				Bold(true)
)

// Check result indicators
const (
	CheckPass = "✓"
	CheckWarn = "⚠"
	CheckFail = "✗"
	Spinner   = "◌"
)

// FormatPass formats a passing check message
func FormatPass(msg string) string {
	return SuccessStyle.Render(CheckPass) + " " + msg
}

// FormatWarn formats a warning message
func FormatWarn(msg string) string {
	return WarningStyle.Render(CheckWarn) + " " + msg
}

// FormatFail formats a failing check message
func FormatFail(msg string) string {
	return ErrorStyle.Render(CheckFail) + " " + msg
}

// FormatMuted formats muted/hint text
func FormatMuted(msg string) string {
	return MutedStyle.Render(msg)
}
