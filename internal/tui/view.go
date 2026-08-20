package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// View renders the current editor state.
func (m Model) View() tea.View {
	var b strings.Builder

	for i := 0; i < m.editor.LineCount(); i++ {
		b.WriteString(m.editor.Line(i))
		b.WriteRune('\n')
	}

	view := tea.NewView(b.String())

	cursor := m.editor.Cursor()

	view.Cursor = &tea.Cursor{
		X:     cursor.Column,
		Y:     cursor.Line,
		Shape: 2,
		Blink: true,
	}

	return view
}
