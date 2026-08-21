package tui

import (
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/i-walk-away/daun/internal/editor"
)

func renderLine(
	text string,
	leftColumn int,
	width int,
) string {
	if width <= 0 {
		return ""
	}

	runes := []rune(text)

	return renderVisibleRunes(
		runes,
		leftColumn,
		width,
		nil,
	)
}

func renderLineSelection(
	text string,
	line int,
	start editor.Position,
	end editor.Position,
	leftColumn int,
	width int,
) string {
	if width <= 0 {
		return ""
	}

	runes := []rune(text)

	selectionStart := 0
	selectionEnd := 0

	switch {
	case line < start.Line || line > end.Line:
		return renderLine(text, leftColumn, width)

	case start.Line == end.Line:
		selectionStart = start.Column
		selectionEnd = end.Column

	case line == start.Line:
		selectionStart = start.Column
		selectionEnd = len(runes)

	case line == end.Line:
		selectionEnd = end.Column
		selectionStart = 0

	default:
		selectionStart = 0
		selectionEnd = len(runes)
	}

	selectionStart = min(max(selectionStart, 0), len(runes))
	selectionEnd = min(max(selectionEnd, selectionStart), len(runes))

	selection := &selectionRange{
		start: selectionStart,
		end:   selectionEnd,
	}

	if len(runes) == 0 {
		return selectionStyle.Render(" ")
	}

	return renderVisibleRunes(
		runes,
		leftColumn,
		width,
		selection,
	)
}

type selectionRange struct {
	start int
	end   int
}

func renderVisibleRunes(
	runes []rune,
	leftColumn int,
	width int,
	selection *selectionRange,
) string {
	rightColumn := leftColumn + width

	var b strings.Builder

	visibleWidth := 0

	for i, r := range runes {
		runeWidth := runewidth.RuneWidth(r)

		if runeWidth == 0 {
			continue
		}

		cellStart := visibleWidth
		cellEnd := visibleWidth + runeWidth

		visibleWidth = cellEnd

		if cellEnd <= leftColumn {
			continue
		}

		if cellStart >= rightColumn {
			break
		}

		// Do not render a wide rune when only part of it would fit.
		if cellStart < leftColumn || cellEnd > rightColumn {
			continue
		}

		selected := selection != nil &&
			i >= selection.start &&
			i < selection.end

		if selected {
			b.WriteString(selectionStyle.Render(string(r)))
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}

func displayWidth(runes []rune) int {
	return runewidth.StringWidth(string(runes))
}
