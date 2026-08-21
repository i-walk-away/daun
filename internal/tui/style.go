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
				Padding(0, 1)

	messageInfoStyle = lipgloss.NewStyle()

	messageSuccessStyle = lipgloss.NewStyle().
				Bold(true)

	messageWarningStyle = lipgloss.NewStyle().
				Bold(true)

	messageErrorStyle = lipgloss.NewStyle().
				Bold(true)

	searchSeparatorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("8"))

	searchMatchStyle = lipgloss.NewStyle().
				Reverse(true)

	searchCurrentMatchStyle = lipgloss.NewStyle().
				Bold(true).
				Reverse(true)

	searchResultStyle = lipgloss.NewStyle()

	searchResultSelectedStyle = lipgloss.NewStyle().
					Reverse(true)

	searchInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))
)
