package editor

// Insert inserts a rune at the current caret position.
//
// Parameters:
//   - r: The Unicode character (rune) to insert into the text.
//
// This function handles the insertion of any valid Unicode character,
// including multi-byte characters like emojis or accented letters.
func (e *Editor) Insert(r rune) {
	e.buffer.Insert(e.cursor.Line, e.cursor.Column, r)
	e.cursor.Column++
}

// Backspace removes the rune immediately before the caret (left side).
//
// If the cursor is at the beginning of a line (column 0), this function
// merges the current line with the previous line.
//
// If the cursor is at the beginning of the entire buffer (line 0, column 0),
// nothing happens.
func (e *Editor) Backspace() {
	// If cursor at the very beginning,
	// do nothing
	if e.cursor.Line == 0 && e.cursor.Column == 0 {
		return
	}

	// If cursor is at the beginning of a line,
	// merge the current line with the previous line.
	if e.cursor.Column == 0 {
		previousLineLength := e.buffer.LineLength(e.cursor.Line - 1)

		e.buffer.JoinLines(e.cursor.Line - 1)

		e.cursor.Line--
		e.cursor.Column = previousLineLength
		return
	}

	e.buffer.Delete(e.cursor.Line, e.cursor.Column)
	e.cursor.Column--
}

// DeleteWordBackwards deletes the text from the caret to the beginning
// of the nearest word to the left.
//
// If the caret is at the beginning of a line, the current line is merged
// with the previous line. If the caret is at the beginning of the buffer,
// nothing happens.
func (e *Editor) DeleteWordBackwards() {
	// If cursor at the very beginning,
	// do nothing
	if e.cursor.Line == 0 && e.cursor.Column == 0 {
		return
	}

	// If cursor is at the beginning of a line,
	// merge the current line with the previous line.
	if e.cursor.Column == 0 {
		previousLineLength := e.buffer.LineLength(e.cursor.Line - 1)

		e.buffer.JoinLines(e.cursor.Line - 1)

		e.cursor.Line--
		e.cursor.Column = previousLineLength
		return
	}

	line := e.buffer.Line(e.cursor.Line)
	start := previousWordBoundary(line, e.cursor.Column)

	e.buffer.DeleteRange(
		e.cursor.Line,
		start,
		e.cursor.Column,
	)

	e.cursor.Column = start
}

// Enter inserts a newline at the current cursor position.
//
// This function splits the current line at the cursor position and
// moves the cursor to the beginning of the new line.
func (e *Editor) Enter() {
	e.buffer.SplitLine(e.cursor.Line, e.cursor.Column)

	e.cursor.Line++
	e.cursor.Column = 0
}
