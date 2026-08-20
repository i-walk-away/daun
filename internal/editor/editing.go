package editor

// Insert inserts a rune at the current caret position.
//
// Parameters:
//   - r: The Unicode character (rune) to insert into the text.
//
// This function handles the insertion of any valid Unicode character,
// including multi-byte characters like emojis or accented letters.
func (e *Editor) Insert(r rune) {
	// Convert current line to rune slice
	var line []rune = e.lineToRuneSlice(e.cursor.Line)

	// Split the line at the cursor position
	leftSide := line[:e.cursor.Column]  // Left side of cursor
	rightSide := line[e.cursor.Column:] // Right side of cursor

	// Create new slice with pre-allocated capacity for performance
	updated := make([]rune, 0, len(line)+1)
	updated = append(updated, leftSide...)  // Add left side
	updated = append(updated, r)            // Insert new character
	updated = append(updated, rightSide...) // Add right side

	// Convert back to string and store in buffer
	e.buffer.SetLine(e.cursor.Line, string(updated))
	// Move cursor to after the inserted character
	e.cursor.Column++
}

// Backspace removes the rune immediately before the caret (left side).
//
// If the cursor is at the beginning of a line (column 0), this function
// should merge the current line with the previous line.
//
// If the cursor is at the beginning of the entire buffer (line 0, column 0),
// nothing should happen.
func (e *Editor) Backspace() {
	// If cursor at the very beginning,
	// do nothing
	if e.cursor.Line == 0 && e.cursor.Column == 0 {
		return
	}

	// If cursor at the beginning of a line (not first line)
	if e.cursor.Column == 0 {
		// Get current line
		currentLine := e.buffer.Line(e.cursor.Line)
		// Merge with previous line
		prevLine := e.buffer.Line(e.cursor.Line - 1)
		e.buffer.SetLine(e.cursor.Line-1, prevLine+currentLine)

		// Remove current line
		e.buffer.DeleteLine(e.cursor.Line)

		// Move cursor to the end of the merged line
		e.cursor.Line--
		e.cursor.Column = e.getLineLength(e.cursor.Line)
		return
	}

	// Normal backspace - remove character before cursor
	line := e.lineToRuneSlice(e.cursor.Line)

	// Create a new slice with capacity for the line minus one character
	// Pre-allocating capacity improves performance by avoiding reallocations
	updated := make([]rune, 0, len(line)-1)

	// Add all characters BEFORE the character being deleted
	updated = append(updated, line[:e.cursor.Column-1]...)

	// Add all characters AFTER the character being deleted
	updated = append(updated, line[e.cursor.Column:]...)

	// Convert back to string and store in buffer
	e.buffer.SetLine(e.cursor.Line, string(updated))

	// Move cursor one position to the left
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
		currentLine := e.buffer.Line(e.cursor.Line)
		prevLine := e.buffer.Line(e.cursor.Line - 1)

		e.buffer.SetLine(e.cursor.Line-1, prevLine+currentLine)
		e.buffer.DeleteLine(e.cursor.Line)

		e.cursor.Line--
		e.cursor.Column = e.getLineLength(e.cursor.Line)
		return
	}

	line := e.lineToRuneSlice(e.cursor.Line)
	start := previousWordBoundary(line, e.cursor.Column)

	updated := make([]rune, 0, len(line)-(e.cursor.Column-start))
	updated = append(updated, line[:start]...)
	updated = append(updated, line[e.cursor.Column:]...)

	e.buffer.SetLine(e.cursor.Line, string(updated))
	e.cursor.Column = start
}

// Enter inserts a newline at the current cursor position.
//
// This function splits the current line at the cursor position and
// moves the cursor to the beginning of the new line.
func (e *Editor) Enter() {
	// Convert current line to rune slice
	var line []rune = e.lineToRuneSlice(e.cursor.Line)

	// Split the line at the cursor position
	leftSide := string(line[:e.cursor.Column])  // Text before cursor
	rightSide := string(line[e.cursor.Column:]) // Text after cursor

	// Replace current line with left side
	e.buffer.SetLine(e.cursor.Line, leftSide)

	// Insert new line with right side after current line
	// Insert at position e.cursor.Line + 1
	e.buffer.InsertLine(e.cursor.Line+1, rightSide)

	// Move cursor to the beginning of the new line
	e.cursor.Line++
	e.cursor.Column = 0
}
