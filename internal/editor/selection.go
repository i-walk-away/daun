package editor

import (
	"strings"
)

// Selection represents an active text selection.
//
// Anchor is the cursor position at which the selection started.
// The other boundary is the editor's current cursor position.
type Selection struct {
	Anchor Position
}

// NewSelection creates a selection starting at anchor.
func NewSelection(anchor Position) *Selection {
	return &Selection{
		Anchor: anchor,
	}
}

// BeginSelection starts a selection at the current cursor position.
//
// If a selection is already active, its anchor is preserved.
func (e *Editor) BeginSelection() {
	if e.selection != nil {
		return
	}

	e.selection = NewSelection(e.cursor)
}

// ClearSelection removes the current selection.
func (e *Editor) ClearSelection() {
	e.selection = nil
}

// SelectAll selects the entire document.
func (e *Editor) SelectAll() {
	if e.buffer.LineCount() == 0 {
		e.selection = nil
		return
	}

	end := Position{
		Line:   e.buffer.LineCount() - 1,
		Column: e.buffer.LineLength(e.buffer.LineCount() - 1),
	}

	if end == (Position{}) {
		e.selection = nil
		return
	}

	e.selection = NewSelection(Position{})
	e.cursor = end
}

// Selection returns the currently selected range.
//
// The returned positions are ordered from the beginning of the buffer
// to the end of the buffer. The boolean is false when no selection is active.
func (e *Editor) Selection() (Position, Position, bool) {
	if e.selection == nil {
		return Position{}, Position{}, false
	}

	start := e.selection.Anchor
	end := e.cursor

	if positionLess(end, start) {
		start, end = end, start
	}

	return start, end, true
}

// SelectedText returns the text currently selected by the editor.
//
// It returns false when there is no non-empty selection.
func (e *Editor) SelectedText() (string, bool) {
	start, end, ok := e.Selection()
	if !ok || start == end {
		return "", false
	}

	return e.buffer.TextRange(
		start.Line,
		start.Column,
		end.Line,
		end.Column,
	), true
}

// CutSelection removes the current selection and returns its text.
//
// It returns false when there is no non-empty selection.
func (e *Editor) CutSelection() (string, bool) {
	text, ok := e.SelectedText()
	if !ok {
		return "", false
	}

	e.DeleteSelection()

	return text, true
}

// positionLess reports whether a comes before b in the buffer.
func positionLess(a, b Position) bool {
	if a.Line != b.Line {
		return a.Line < b.Line
	}

	return a.Column < b.Column
}

// normalizePaste converts common platform line endings to \n.
func normalizePaste(text string) string {
	return strings.ReplaceAll(text, "\r\n", "\n")
}
