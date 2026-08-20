package buffer

// Line represents a single line of text.
type Line struct {
	content gapBuffer
}

// newLine creates a line containing text.
func newLine(text string) Line {
	return Line{
		content: newGapBuffer(text),
	}
}

// Insert inserts a rune at the current editing position.
func (l *Line) Insert(r rune) {
	l.content.Insert(r)
}

// Delete removes the rune immediately before the current editing position.
func (l *Line) Delete() {
	l.content.Delete()
}

func (l *Line) DeleteRange(start, end int) {
	l.content.DeleteRange(start, end)
}

// SetPosition moves the editing position to the given character offset.
func (l *Line) SetPosition(position int) {
	l.content.MoveGap(position)
}

// Len returns the number of runes in the line.
func (l *Line) Len() int {
	return l.content.Len()
}

// String returns the complete text of the line.
func (l *Line) String() string {
	return l.content.String()
}
