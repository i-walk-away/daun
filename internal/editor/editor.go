package editor

import buffer2 "github.com/i-walk-away/daun/internal/buffer"

// Editor represents the main text editor instance.
// It holds the current state including the text buffer and cursor position.
type Editor struct {
	buffer    *buffer2.Buffer
	cursor    Position
	selection *Selection
}

// Position represents the cursor location in the editor.
// Both Line and Column are zero-based indices.
type Position struct {
	Line   int // Current line number (0 = first line)
	Column int // Current column position (0 = first character)
}

// NewEditor creates a new editor from an existing buffer.
// It takes a pointer to the buffer so changes are reflected everywhere.
func NewEditor(buffer *buffer2.Buffer) *Editor {
	return &Editor{
		buffer: buffer,
		cursor: Position{0, 0},
	}
}

// Cursor returns the current cursor position.
func (e *Editor) Cursor() Position {
	return e.cursor
}

// Line returns the text of the line at index.
func (e *Editor) Line(index int) string {
	return e.buffer.Line(index)
}

// LineCount returns the number of lines in the editor.
func (e *Editor) LineCount() int {
	return e.buffer.LineCount()
}

// getLineLength returns the number of Unicode characters (runes) in a line.
func (e *Editor) getLineLength(lineNum int) int {
	return e.buffer.LineLength(lineNum)
}
