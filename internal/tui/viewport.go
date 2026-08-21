package tui

import (
	"github.com/mattn/go-runewidth"

	"github.com/i-walk-away/daun/internal/editor"
)

// Viewport controls which part of the document is visible on screen.
type Viewport struct {
	topLine    int
	leftColumn int
	width      int
	height     int
}

// SetSize updates the viewport dimensions.
func (v *Viewport) SetSize(width, height int) {
	v.width = max(width, 0)
	v.height = max(height, 0)

	v.clamp()
}

// TopLine returns the first visible document line.
func (v Viewport) TopLine() int {
	return v.topLine
}

// LeftColumn returns the first visible terminal column.
func (v Viewport) LeftColumn() int {
	return v.leftColumn
}

// Width returns the viewport width in terminal cells.
func (v Viewport) Width() int {
	return v.width
}

// Height returns the viewport height in terminal rows.
func (v Viewport) Height() int {
	return v.height
}

// ContentWidth returns the number of terminal cells available for text.
func (v Viewport) ContentWidth(lineNumberWidth int) int {
	return max(v.width-lineNumberWidth-3, 0)
}

// EnsureCursorVisible scrolls the viewport so that the cursor remains visible.
func (v *Viewport) EnsureCursorVisible(
	cursor editor.Position,
	line string,
	lineCount int,
	lineNumberWidth int,
) {
	if v.width <= 0 || v.height <= 0 {
		return
	}

	contentWidth := v.ContentWidth(lineNumberWidth)

	if contentWidth <= 0 {
		return
	}

	if cursor.Line < v.topLine {
		v.topLine = cursor.Line
	}

	if cursor.Line >= v.topLine+v.height {
		v.topLine = cursor.Line - v.height + 1
	}

	runes := []rune(line)
	cursorColumn := min(cursor.Column, len(runes))
	cursorX := runewidth.StringWidth(string(runes[:cursorColumn]))

	if cursorX < v.leftColumn {
		v.leftColumn = cursorX
	}

	if cursorX >= v.leftColumn+contentWidth {
		v.leftColumn = cursorX - contentWidth + 1
	}

	v.clamp()

	maxTopLine := max(lineCount-v.height, 0)

	if v.topLine > maxTopLine {
		v.topLine = maxTopLine
	}
}

// VisibleLines returns the first and last visible document line.
func (v Viewport) VisibleLines(lineCount int) (int, int) {
	if lineCount <= 0 || v.height <= 0 {
		return 0, 0
	}

	start := min(v.topLine, lineCount)
	end := min(start+v.height, lineCount)

	return start, end
}

func (v *Viewport) clamp() {
	v.topLine = max(v.topLine, 0)
	v.leftColumn = max(v.leftColumn, 0)
}

// ScrollUp moves the viewport up by the specified number of lines.
func (v *Viewport) ScrollUp(lines int) {
	v.topLine = max(v.topLine-lines, 0)
}

// ScrollDown moves the viewport down by the specified number of lines.
func (v *Viewport) ScrollDown(lines, lineCount int) {
	maxTopLine := max(lineCount-v.height, 0)
	v.topLine = min(v.topLine+lines, maxTopLine)
}
