package editor

// Buffer holds the text content of the editor.
// The text is stored as a slice of strings, where each string represents
// a single line.
type Buffer struct {
	lines []string
}

// NewBuffer creates a new buffer with one empty line.
// This ensures the buffer is never empty and always has at least one line.
func NewBuffer() *Buffer {
	return &Buffer{
		lines: []string{""}, // Start with one empty line
	}
}

// Line returns the text of the line at index.
func (b *Buffer) Line(index int) string {
	return b.lines[index]
}

// SetLine replaces the text of the line at index.
func (b *Buffer) SetLine(index int, line string) {
	b.lines[index] = line
}

// InsertLine inserts a new line at index.
// The existing line at index and all subsequent lines are shifted down.
func (b *Buffer) InsertLine(index int, line string) {
	b.lines = append(b.lines, "")
	copy(b.lines[index+1:], b.lines[index:])
	b.lines[index] = line
}

// DeleteLine removes the line at index.
// All subsequent lines are shifted up.
func (b *Buffer) DeleteLine(index int) {
	copy(b.lines[index:], b.lines[index+1:])
	b.lines = b.lines[:len(b.lines)-1]
}

// LineCount returns the number of lines in the buffer.
func (b *Buffer) LineCount() int {
	return len(b.lines)
}

// LineLength returns the number of runes in the line at index.
func (b *Buffer) LineLength(index int) int {
	return len([]rune(b.lines[index]))
}
