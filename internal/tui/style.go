package tui

import "github.com/charmbracelet/lipgloss"

var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())

	statusBarStyle = lipgloss.NewStyle().
			Reverse(true)

	lineNumberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	currentLineNumberStyle = lipgloss.NewStyle().
				Bold(true)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	selectionStyle = lipgloss.NewStyle().
			Reverse(true)
)
