package editor

// MoveLeft moves the caret one character to the left.
//
// If the cursor is at the beginning of a line (column 0), moving left
// wraps to the end of the previous line.
//
// If the cursor is at the very beginning of the buffer (line 0, column 0),
// the cursor does not move.
func (e *Editor) MoveLeft() {
	// If caret already at the very beginning of the buffer,
	// do nothing
	if e.cursor.Column == 0 && e.cursor.Line == 0 {
		return
	}

	// If caret is at the beginning of a line,
	// wrap to the end of the previous line
	if e.cursor.Column == 0 {
		e.cursor.Line--
		// Move caret to the end of the previous line
		e.cursor.Column = e.getLineLength(e.cursor.Line)
		return
	}

	// Otherwise, simply move left by 1 column
	e.cursor.Column--
}

// MoveRight moves the caret one character to the right.
//
// If the cursor is at the end of a line, moving right wraps to the
// beginning of the next line.
//
// If the cursor is at the very end of the buffer (last line, end of line),
// the cursor does not move.
func (e *Editor) MoveRight() {
	// If caret is at the very end of the buffer,
	// do nothing
	if e.cursor.Line == e.buffer.LineCount()-1 &&
		e.cursor.Column == e.getLineLength(e.cursor.Line) {
		return
	}

	// If caret is at the end of a line,
	// wrap to the beginning of the next line
	if e.cursor.Column == e.getLineLength(e.cursor.Line) {
		e.cursor.Line++
		e.cursor.Column = 0
		return
	}

	// Otherwise, simply move right by 1 column
	e.cursor.Column++
}

// MoveUp moves the caret up by one line.
func (e *Editor) MoveUp() {
	// If already at the first line,
	// move cursor to the beginning of the line
	if e.cursor.Line == 0 {
		e.cursor.Column = 0
		return
	}

	// If the target line length is shorter than the current column,
	// place the cursor at the end of the target line
	targetLineLength := e.getLineLength(e.cursor.Line - 1)
	if targetLineLength < e.cursor.Column {
		e.cursor.Column = targetLineLength
	}

	// cursor.Column is already preserved
	// if the `if` statement above is not executed
	e.cursor.Line--
}

// MoveDown moves the caret down by one line.
func (e *Editor) MoveDown() {
	// If already at the last line,
	// move cursor to the end of the line
	// LineCount returns the number of lines, while line indexes start at 0,
	// so the last line has index LineCount() - 1.
	if e.cursor.Line == e.buffer.LineCount()-1 {
		e.cursor.Column = e.getLineLength(e.cursor.Line)
		return
	}

	// If the target line length is shorter than the current column,
	// place the cursor at the end of the target line
	targetLineLength := e.getLineLength(e.cursor.Line + 1)
	if targetLineLength < e.cursor.Column {
		e.cursor.Column = targetLineLength
	}

	// cursor.Column is preserved if the `if` statement above is not executed
	e.cursor.Line++
}
