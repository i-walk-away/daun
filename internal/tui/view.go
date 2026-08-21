package tui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/i-walk-away/daun/internal/editor"
)

// View renders the editor.
func (m Model) View() tea.View {
	var builder strings.Builder

	start, end, hasSelection := m.editor.Selection()

	lineNumberWidth := lineNumberWidth(
		m.editor.LineCount(),
	)

	visibleStart, visibleEnd := m.viewport.VisibleLines(
		m.editor.LineCount(),
	)

	contentHeight := m.viewport.Height()

	for line := visibleStart; line < visibleEnd; line++ {
		rawNumber := strconv.Itoa(line + 1)

		padding := lineNumberWidth -
			len([]rune(rawNumber))

		if padding > 0 {
			builder.WriteString(
				strings.Repeat(" ", padding),
			)
		}

		if line == m.editor.Cursor().Line {
			builder.WriteString(
				currentLineNumberStyle.Render(
					rawNumber,
				),
			)
		} else {
			builder.WriteString(
				lineNumberStyle.Render(
					rawNumber,
				),
			)
		}

		builder.WriteString(
			separatorStyle.Render(" │ "),
		)

		text := m.editor.Line(line)

		contentWidth := m.viewport.ContentWidth(
			lineNumberWidth,
		)

		switch {
		case m.search.Mode == searchModeText:
			var current *editor.Match

			if m.search.text.Found {
				match := m.search.text.Current
				current = &match
			}

			text = renderLineSearch(
				text,
				line,
				m.search.Input.Value(),
				current,
				m.viewport.LeftColumn(),
				contentWidth,
			)

		case hasSelection:
			text = renderLineSelection(
				text,
				line,
				start,
				end,
				m.viewport.LeftColumn(),
				contentWidth,
			)

		default:
			text = renderLine(
				text,
				m.viewport.LeftColumn(),
				contentWidth,
			)
		}

		builder.WriteString(text)
		builder.WriteRune('\n')
	}

	renderedLines := visibleEnd - visibleStart

	for renderedLines < contentHeight {
		builder.WriteRune('\n')
		renderedLines++
	}

	// Search occupies the area directly above the message panel.
	if m.search.Mode != 0 {
		builder.WriteString(
			m.renderSearchPanel(),
		)
		builder.WriteRune('\n')
	}

	// Messages remain visible while search is open.
	if m.message != nil {
		builder.WriteString(
			m.renderMessage(),
		)
		builder.WriteRune('\n')
	}

	builder.WriteString(
		m.renderStatusBar(),
	)

	view := tea.NewView(
		builder.String(),
	)

	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	if m.search.Mode != 0 {
		m.setSearchCursor(
			&view,
			contentHeight,
		)
	} else {
		m.setCursor(
			&view,
			lineNumberWidth,
		)
	}

	return view
}

func (m Model) renderMessage() string {
	if m.message == nil {
		return ""
	}

	style := messageInfoStyle

	switch m.message.Level {
	case MessageSuccess:
		style = messageSuccessStyle

	case MessageWarning:
		style = messageWarningStyle

	case MessageError:
		style = messageErrorStyle
	}

	return messagePanelStyle.
		Width(max(m.width, 1)).
		Render(
			style.Render(m.message.Text),
		)
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
		return statusBarStyle.Render(
			left + "  " + right,
		)
	}

	padding := m.width -
		len([]rune(left)) -
		len([]rune(right))

	if padding < 1 {
		return statusBarStyle.Render(left)
	}

	return statusBarStyle.Render(
		left + strings.Repeat(" ", padding) + right,
	)
}

func (m Model) setCursor(
	view *tea.View,
	lineNumberWidth int,
) {
	cursor := m.editor.Cursor()
	line := []rune(
		m.editor.Line(cursor.Line),
	)

	cursorColumn := min(
		cursor.Column,
		len(line),
	)

	cursorX := displayWidth(
		line[:cursorColumn],
	)

	cursorX -= m.viewport.LeftColumn()
	cursorX += lineNumberWidth + 3

	cursorY := cursor.Line -
		m.viewport.TopLine()

	view.Cursor = &tea.Cursor{
		X: max(
			cursorX,
			lineNumberWidth+3,
		),
		Y:     max(cursorY, 0),
		Shape: 2,
		Blink: true,
	}
}

func (m Model) setSearchCursor(
	view *tea.View,
	contentHeight int,
) {
	cursor := m.search.Input.Cursor()

	if cursor == nil {
		return
	}

	cursor.Y = contentHeight + 1

	view.Cursor = cursor
}

func lineNumberWidth(lineCount int) int {
	return len(
		strconv.Itoa(
			max(lineCount, 1),
		),
	)
}
