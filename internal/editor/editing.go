package editor

// Insert inserts a rune at the current caret position.
//
// If text is selected, the selection is replaced by the inserted rune.
func (e *Editor) Insert(r rune) {
	e.InsertText(string(r))
}

// InsertText inserts text at the current caret position.
//
// If text is selected, the selection is replaced by the inserted text.
// Newline characters create new lines.
func (e *Editor) InsertText(text string) {
	text = normalizePaste(text)

	start := e.cursor
	end := e.cursor

	if selectionStart, selectionEnd, ok := e.Selection(); ok &&
		selectionStart != selectionEnd {
		start = selectionStart
		end = selectionEnd
	}

	afterCursor := advancePosition(start, text)

	e.applyEdit(
		start,
		end,
		text,
		afterCursor,
		nil,
	)
}

// Backspace removes the rune immediately before the caret.
//
// If text is selected, the entire selection is deleted instead.
func (e *Editor) Backspace() {
	if start, end, ok := e.Selection(); ok && start != end {
		e.deleteRange(start, end)
		return
	}

	if e.cursor.Line == 0 && e.cursor.Column == 0 {
		return
	}

	if e.cursor.Column == 0 {
		start := Position{
			Line:   e.cursor.Line - 1,
			Column: e.buffer.LineLength(e.cursor.Line - 1),
		}

		end := e.cursor

		e.applyEdit(
			start,
			end,
			"",
			start,
			nil,
		)

		return
	}

	start := Position{
		Line:   e.cursor.Line,
		Column: e.cursor.Column - 1,
	}

	end := e.cursor

	e.applyEdit(
		start,
		end,
		"",
		start,
		nil,
	)
}

// DeleteSelection removes all text inside the current selection.
//
// After deletion, the cursor is placed at the beginning of the selection.
func (e *Editor) DeleteSelection() {
	start, end, ok := e.Selection()
	if !ok || start == end {
		return
	}

	e.deleteRange(start, end)
}

// Enter inserts a newline at the current cursor position.
//
// If text is selected, the selection is replaced by the newline.
func (e *Editor) Enter() {
	start := e.cursor
	end := e.cursor

	if selectionStart, selectionEnd, ok := e.Selection(); ok &&
		selectionStart != selectionEnd {
		start = selectionStart
		end = selectionEnd
	}

	afterCursor := Position{
		Line:   start.Line + 1,
		Column: 0,
	}

	e.applyEdit(
		start,
		end,
		"\n",
		afterCursor,
		nil,
	)
}

// DeleteWordBackwards deletes text from the caret to the beginning
// of the nearest word to the left.
//
// If the caret is at the beginning of a line, the current line is merged
// with the previous line.
func (e *Editor) DeleteWordBackwards() {
	if start, end, ok := e.Selection(); ok && start != end {
		e.deleteRange(start, end)
		return
	}

	if e.cursor.Line == 0 && e.cursor.Column == 0 {
		return
	}

	if e.cursor.Column == 0 {
		start := Position{
			Line:   e.cursor.Line - 1,
			Column: e.buffer.LineLength(e.cursor.Line - 1),
		}

		e.applyEdit(
			start,
			e.cursor,
			"",
			start,
			nil,
		)

		return
	}

	line := e.buffer.Line(e.cursor.Line)
	column := e.cursor.Column

	runes := []rune(line)

	for column > 0 && isWhitespace(runes[column-1]) {
		column--
	}

	for column > 0 && !isWhitespace(runes[column-1]) {
		column--
	}

	start := Position{
		Line:   e.cursor.Line,
		Column: column,
	}

	e.applyEdit(
		start,
		e.cursor,
		"",
		start,
		nil,
	)
}

func (e *Editor) deleteRange(start, end Position) {
	e.applyEdit(
		start,
		end,
		"",
		start,
		nil,
	)
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}
