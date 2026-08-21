package tui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// View renders the editor.
func (m Model) View() tea.View {
	var b strings.Builder

	start, end, hasSelection := m.editor.Selection()

	lineNumberWidth := lineNumberWidth(m.editor.LineCount())
	visibleStart, visibleEnd := m.viewport.VisibleLines(m.editor.LineCount())
	contentHeight := m.viewport.Height()

	for line := visibleStart; line < visibleEnd; line++ {
		rawNumber := strconv.Itoa(line + 1)

		b.WriteString(
			strings.Repeat(
				" ",
				lineNumberWidth-len([]rune(rawNumber)),
			),
		)

		if line == m.editor.Cursor().Line {
			b.WriteString(currentLineNumberStyle.Render(rawNumber))
		} else {
			b.WriteString(lineNumberStyle.Render(rawNumber))
		}

		b.WriteString(separatorStyle.Render(" │ "))

		text := m.editor.Line(line)
		contentWidth := m.viewport.ContentWidth(lineNumberWidth)

		if hasSelection {
			text = renderLineSelection(
				text,
				line,
				start,
				end,
				m.viewport.LeftColumn(),
				contentWidth,
			)
		} else {
			text = renderLine(
				text,
				m.viewport.LeftColumn(),
				contentWidth,
			)
		}

		b.WriteString(text)
		b.WriteRune('\n')
	}

	renderedLines := visibleEnd - visibleStart

	for renderedLines < contentHeight {
		b.WriteRune('\n')
		renderedLines++
	}

	if m.message != nil {
		b.WriteString(m.renderMessage())
		b.WriteRune('\n')
	}

	b.WriteString(m.renderStatusBar())

	view := tea.NewView(b.String())
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	m.setCursor(&view, lineNumberWidth)

	return view
}

func (m Model) renderMessage() string {
	if m.message == nil {
		return ""
	}

	var style = messageInfoStyle

	switch m.message.Level {
	case MessageSuccess:
		style = messageSuccessStyle

	case MessageWarning:
		style = messageWarningStyle

	case MessageError:
		style = messageErrorStyle
	}

	text := style.Render(m.message.Text)

	return messagePanelStyle.Width(max(m.width-2, 0)).Render(text)
}

func (m Model) renderStatusBar() string {
	cursor := m.editor.Cursor()

	fileName := "[No Name]"

	if path := m.editor.FilePath(); path != "" {
		fileName = filepath.Base(path)
	}

	modifiedMarker := ""

	if m.editor.Modified() {
		modifiedMarker = " [+]"
	}

	left := fileName + modifiedMarker

	right := fmt.Sprintf(
		"Ln %d, Col %d   UTF-8",
		cursor.Line+1,
		cursor.Column+1,
	)

	if m.width <= 0 {
		return statusBarStyle.Render(left + "  " + right)
	}

	padding := m.width - len([]rune(left)) - len([]rune(right))

	if padding < 1 {
		return statusBarStyle.Render(left)
	}

	return statusBarStyle.Render(
		left + strings.Repeat(" ", padding) + right,
	)
}

func (m Model) setCursor(view *tea.View, lineNumberWidth int) {
	cursor := m.editor.Cursor()
	line := []rune(m.editor.Line(cursor.Line))

	cursorColumn := min(cursor.Column, len(line))
	cursorX := displayWidth(line[:cursorColumn])
	cursorX -= m.viewport.LeftColumn()
	cursorX += lineNumberWidth + 3

	cursorY := cursor.Line - m.viewport.TopLine()

	view.Cursor = &tea.Cursor{
		X:     max(cursorX, lineNumberWidth+3),
		Y:     max(cursorY, 0),
		Shape: 2,
		Blink: true,
	}
}

func lineNumberWidth(lineCount int) int {
	return len(strconv.Itoa(max(lineCount, 1)))
}
