package editor

// Insert inserts a rune at the current caret position.
//
// If text is selected, the selection is deleted first and the rune is
// inserted at the beginning of the selection.
func (e *Editor) Insert(r rune) {
	if start, end, ok := e.Selection(); ok && start != end {
		e.DeleteSelection()
	}

	e.buffer.Insert(e.cursor.Line, e.cursor.Column, r)
	e.cursor.Column++

	e.selection = nil
}

// Backspace removes the rune immediately before the caret.
//
// If text is selected, the entire selection is deleted instead.
func (e *Editor) Backspace() {
	if start, end, ok := e.Selection(); ok && start != end {
		e.DeleteSelection()
		return
	}

	if e.cursor.Line == 0 && e.cursor.Column == 0 {
		return
	}

	if e.cursor.Column == 0 {
		e.buffer.JoinLines(e.cursor.Line - 1)

		e.cursor.Line--
		e.cursor.Column = e.buffer.LineLength(e.cursor.Line)
		return
	}

	e.buffer.Delete(e.cursor.Line, e.cursor.Column)
	e.cursor.Column--
}

// DeleteSelection removes all text inside the current selection.
//
// After deletion, the cursor is placed at the beginning of the selection.
func (e *Editor) DeleteSelection() {
	start, end, ok := e.Selection()
	if !ok || start == end {
		return
	}

	e.buffer.DeleteRangeLines(
		start.Line,
		start.Column,
		end.Line,
		end.Column,
	)

	e.cursor = start
	e.selection = nil
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
