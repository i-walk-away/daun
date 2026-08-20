package editor

import "unicode"

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

// MoveWordLeft moves the caret to the beginning of the nearest word
// to the left.
//
// Whitespace is skipped and is never treated as a word. If the caret is
// inside a word, it moves to the beginning of that word. If the caret is
// separated from the previous word by whitespace, the whitespace is
// skipped first and the caret moves to the beginning of that word.
//
// If the caret is at the beginning of a line, it wraps to the end of the
// previous line.
func (e *Editor) MoveWordLeft() {
	if e.cursor.Line == 0 && e.cursor.Column == 0 {
		return
	}

	if e.cursor.Column == 0 {
		e.cursor.Line--
		e.cursor.Column = e.buffer.LineLength(e.cursor.Line)
		return
	}

	e.cursor.Column = previousWordBoundary(
		e.buffer.Line(e.cursor.Line),
		e.cursor.Column,
	)
}

// MoveWordRight moves the caret to the end of the nearest word
// to the right.
//
// If the caret is inside a word, it moves to the end of that word.
// If the caret is on whitespace, the whitespace is skipped and the caret
// moves to the end of the next word.
//
// If there is no word to the right, the caret moves to the end of the line.
//
// If the caret is at the end of a line, it wraps to the beginning of the
// next line.
func (e *Editor) MoveWordRight() {
	if e.cursor.Line == e.buffer.LineCount()-1 &&
		e.cursor.Column == e.buffer.LineLength(e.cursor.Line) {
		return
	}

	if e.cursor.Column == e.buffer.LineLength(e.cursor.Line) {
		e.cursor.Line++
		e.cursor.Column = 0
		return
	}

	e.cursor.Column = nextWordBoundary(
		e.buffer.Line(e.cursor.Line),
		e.cursor.Column,
	)
}

func previousWordBoundary(line string, column int) int {
	runes := []rune(line)

	// Skip whitespace to the left of the cursor.
	for column > 0 && unicode.IsSpace(runes[column-1]) {
		column--
	}

	// Move to the beginning of the previous word.
	for column > 0 && !unicode.IsSpace(runes[column-1]) {
		column--
	}

	return column
}

func nextWordBoundary(line string, column int) int {
	runes := []rune(line)

	// Skip whitespace to the right of the cursor.
	for column < len(runes) && unicode.IsSpace(runes[column]) {
		column++
	}

	// Move to the end of the next word.
	for column < len(runes) && !unicode.IsSpace(runes[column]) {
		column++
	}

	return column
}
