package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// ══════════════════════════════════════════════════════════════
// Colors
// ══════════════════════════════════════════════════════════════

var (
	primaryColor   = lipgloss.Color("#FF0040")
	secondaryColor = lipgloss.Color("#0080FF")
	darkColor      = lipgloss.Color("#1a1b26")
	lightColor     = lipgloss.Color("#c0caf5")
	warningColor   = lipgloss.Color("#ff9e64")
	errorColor     = lipgloss.Color("#f7768e")
	successColor   = lipgloss.Color("#9ece6a")
	infoColor      = lipgloss.Color("#7aa2f7")
)

// ══════════════════════════════════════════════════════════════
// Styles
// ══════════════════════════════════════════════════════════════

var (
	bannerStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Background(darkColor).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Bold(true).
			Padding(0, 1)

	normalStyle = lipgloss.NewStyle().
			Foreground(lightColor)

	disabledStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565f89")).
			Faint(true)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(infoColor)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	statusBarStyle = lipgloss.NewStyle().
			Background(darkColor).
			Foreground(lightColor).
			Padding(0, 1)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(lightColor)
)
