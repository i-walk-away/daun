package editor

// SetCursor moves the cursor to the specified document position.
//
// The supplied position is clamped to the valid document boundaries.
// Setting the cursor also clears the current selection.
func (e *Editor) SetCursor(position Position) {
	e.cursor = e.clampPosition(position)
	e.selection = nil
}

func (e *Editor) clampPosition(position Position) Position {
	lineCount := e.LineCount()

	if lineCount == 0 {
		return Position{}
	}

	position.Line = min(
		max(position.Line, 0),
		lineCount-1,
	)

	lineLength := e.getLineLength(position.Line)

	position.Column = min(
		max(position.Column, 0),
		lineLength,
	)

	return position
}
