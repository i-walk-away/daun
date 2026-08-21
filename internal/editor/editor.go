package editor

import "github.com/i-walk-away/daun/internal/buffer"

// Editor represents the main text editor instance.
//
// Editor owns the current document state, cursor, selection, file path,
// and editing history.
type Editor struct {
	buffer *buffer.Buffer

	cursor    Position
	selection *Selection

	filePath string

	revision      uint64
	savedRevision uint64

	undoStack []historyEntry
	redoStack []historyEntry
}

// Position represents the cursor location in the document.
//
// Both Line and Column are zero-based indices.
type Position struct {
	Line   int
	Column int
}

// NewEditor creates a new editor from a buffer without associating it
// with a file.
func NewEditor(buffer *buffer.Buffer) *Editor {
	return &Editor{
		buffer: buffer,
	}
}

// NewFileEditor creates a new editor associated with path.
func NewFileEditor(buffer *buffer.Buffer, path string) *Editor {
	return &Editor{
		buffer:   buffer,
		filePath: path,
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

// FilePath returns the path of the file associated with the editor.
//
// It returns an empty string for an unsaved document.
func (e *Editor) FilePath() string {
	return e.filePath
}

// Modified reports whether the document contains changes that have not
// been saved.
func (e *Editor) Modified() bool {
	return e.revision != e.savedRevision
}

// Close releases resources owned by the editor's document.
func (e *Editor) Close() error {
	if e == nil || e.buffer == nil {
		return nil
	}

	return e.buffer.Close()
}

// getLineLength returns the number of Unicode characters in a line.
func (e *Editor) getLineLength(line int) int {
	return e.buffer.LineLength(line)
}
