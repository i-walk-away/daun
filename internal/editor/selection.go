package editor

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

// positionLess returns true if `a` comes before `b` in the buffer.
func positionLess(a, b Position) bool {
	if a.Line != b.Line {
		return a.Line < b.Line
	}

	return a.Column < b.Column
}
