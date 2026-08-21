package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"

	"github.com/i-walk-away/daun/internal/editor"
)

var selectionStyle = lipgloss.NewStyle().Reverse(true)

// View renders the current editor state.
func (m Model) View() tea.View {
	var b strings.Builder

	start, end, hasSelection := m.editor.Selection()

	for line := 0; line < m.editor.LineCount(); line++ {
		text := m.editor.Line(line)

		if hasSelection {
			text = renderLineSelection(text, line, start, end)
		}

		b.WriteString(text)
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

// renderLineSelection renders a single line with the selected characters styled.
//
// Selection boundaries are represented as [start, end), where start is
// inclusive and end is exclusive.
func renderLineSelection(
	text string,
	line int,
	start editor.Position,
	end editor.Position,
) string {
	runes := []rune(text)

	// This line is outside the selection.
	if line < start.Line || line > end.Line {
		return text
	}

	// Empty lines inside a multi-line selection need a visible
	// selected cell so that the selection remains visually continuous.
	if len(runes) == 0 && start.Line != end.Line {
		return selectionStyle.Render(" ")
	}

	lineStart := 0
	lineEnd := len(runes)

	switch {
	case start.Line == end.Line:
		lineStart = start.Column
		lineEnd = end.Column

	case line == start.Line:
		lineStart = start.Column

	case line == end.Line:
		lineEnd = end.Column
	}

	// Cursor positions can refer to a column beyond the current line.
	// Clamp both boundaries before slicing the rune slice.
	lineStart = max(0, min(lineStart, len(runes)))
	lineEnd = max(0, min(lineEnd, len(runes)))

	// Nothing to select on this line.
	if lineStart >= lineEnd {
		return text
	}

	var b strings.Builder

	b.WriteString(string(runes[:lineStart]))
	b.WriteString(selectionStyle.Render(string(runes[lineStart:lineEnd])))
	b.WriteString(string(runes[lineEnd:]))

	return b.String()
}
