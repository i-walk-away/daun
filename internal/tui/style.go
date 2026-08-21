package tui

import "github.com/charmbracelet/lipgloss"

var (
	lineNumberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	currentLineNumberStyle = lipgloss.NewStyle().
				Bold(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	selectionStyle = lipgloss.NewStyle().
			Reverse(true)

	statusBarStyle = lipgloss.NewStyle().
			Reverse(true)

	messagePanelStyle = lipgloss.NewStyle().
				BorderTop(true).
				Padding(0, 1)

	messageInfoStyle = lipgloss.NewStyle()

	messageSuccessStyle = lipgloss.NewStyle().
				Bold(true)

	messageWarningStyle = lipgloss.NewStyle().
				Bold(true)

	messageErrorStyle = lipgloss.NewStyle().
				Bold(true)
)
