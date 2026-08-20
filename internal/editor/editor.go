package editor

// Editor represents the main text editor instance.
// It holds the current state including the text buffer and cursor position.
type Editor struct {
	buffer Buffer
	cursor Position
}

// Position represents the cursor location in the editor.
// Both Line and Column are zero-based indices.
type Position struct {
	Line   int // Current line number (0 = first line)
	Column int // Current column position (0 = first character)
}

// Buffer holds the text content of the editor.
// The text is stored as a slice of strings, where each string represents
// a single line. This is a simple but effective approach for small to medium files.
type Buffer struct {
	lines []string // Each element is one line of text
}

// Insert inserts a rune at the current caret position.
//
// Parameters:
//   - r: The Unicode character (rune) to insert into the text.
//
// This function handles the insertion of any valid Unicode character,
// including multi-byte characters like emojis or accented letters.
func (e *Editor) Insert(r rune) {
	// Convert current line to rune slice
	line := []rune(e.buffer.lines[e.cursor.Line])

	// Split the line at the cursor position
	leftSide := line[:e.cursor.Column]  // Left side of cursor
	rightSide := line[e.cursor.Column:] // Right side of cursor

	// Create new slice with pre-allocated capacity for performance
	updated := make([]rune, 0, len(line)+1)
	updated = append(updated, leftSide...)  // Add left side
	updated = append(updated, r)            // Insert new character
	updated = append(updated, rightSide...) // Add right side

	// Convert back to string and store in buffer
	e.buffer.lines[e.cursor.Line] = string(updated)
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
//
// Example:
//   - Before: line = "Hello", cursor at column 3 (after 'e')
//   - Backspace() removes 'e'
//   - After:  "Hllo", cursor at column 2
func (e *Editor) Backspace() {
	// TODO: Implement
}

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
	if e.cursor.Line == len(e.buffer.lines)-1 &&
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
//
// The new column position is calculated as follows:
//   - Try to keep the same column position as before
//   - If the target line is shorter than the current column,
//     place the cursor at the end of the target line
//   - If already at the first line (line 0), move cursor to the beginning of the line
func (e *Editor) MoveUp() {
	// If already at the first line,
	// move cursor to the beginning of the line
	if e.cursor.Line == 0 {
		e.cursor.Column = 0
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
//
// The new column position is calculated as follows:
//   - Try to keep the same column position as before
//   - If the target line is shorter than the current column,
//     place the cursor at the end of the target line
//   - If already at the last line, move to the end of the line
func (e *Editor) MoveDown() {
	// If already at the last line,
	// move to the end of the line
	if e.cursor.Line == len(e.buffer.lines)-1 {
		e.cursor.Column = e.getLineLength(e.cursor.Line)
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

// Enter inserts a newline at the current cursor position.
//
// This function splits the current line at the cursor position and
// moves the cursor to the beginning of the new line.
//
// Process:
//  1. Split current line into "before" (left of cursor) and "after" (right of cursor)
//  2. Replace current line with the "before" part
//  3. Insert a new line containing the "after" part
//  4. Move cursor to the beginning of the new line (column 0)
func (e *Editor) Enter() {
	// Convert current line to rune slice
	line := []rune(e.buffer.lines[e.cursor.Line])

	// Split the line at the cursor position
	leftSide := string(line[:e.cursor.Column])  // Text before cursor
	rightSide := string(line[e.cursor.Column:]) // Text after cursor

	// Replace current line with left side
	e.buffer.lines[e.cursor.Line] = leftSide

	// Insert new line with right side after current line
	// Insert at position e.cursor.Line + 1
	e.buffer.lines = append(e.buffer.lines[:e.cursor.Line+1],
		append([]string{rightSide}, e.buffer.lines[e.cursor.Line+1:]...)...)

	// Move cursor to the beginning of the new line
	e.cursor.Line++
	e.cursor.Column = 0
}

// getLineLength returns the number of Unicode characters (runes) in a line
func (e *Editor) getLineLength(lineNum int) int {
	if lineNum < 0 || lineNum >= len(e.buffer.lines) {
		return 0
	}

	return len([]rune(e.buffer.lines[lineNum]))
}
