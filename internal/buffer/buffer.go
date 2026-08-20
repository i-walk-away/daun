package buffer

// Buffer holds the text content of the editor.
// The text is stored as a slice of lines.
type Buffer struct {
	lines []Line
}

// NewBuffer creates a new buffer containing one empty line.
func NewBuffer() *Buffer {
	return &Buffer{
		lines: []Line{{}},
	}
}

// Line returns the text of the line at index.
func (b *Buffer) Line(index int) string {
	return b.lines[index].String()
}

// LineLength returns the number of runes in the line at index.
func (b *Buffer) LineLength(index int) int {
	return b.lines[index].Len()
}

// LineCount returns the number of lines in the buffer.
func (b *Buffer) LineCount() int {
	return len(b.lines)
}

// Insert inserts a rune at the specified position.
func (b *Buffer) Insert(line, column int, r rune) {
	l := &b.lines[line]

	l.SetPosition(column)
	l.Insert(r)
}

// Delete removes the rune immediately before the specified position.
func (b *Buffer) Delete(line, column int) {
	l := &b.lines[line]

	l.SetPosition(column)
	l.Delete()
}

// InsertLine inserts an empty line at the specified index.
func (b *Buffer) InsertLine(index int) {
	b.lines = append(b.lines, Line{})
	copy(b.lines[index+1:], b.lines[index:])
	b.lines[index] = Line{}
}

// DeleteLine removes the line at the specified index.
func (b *Buffer) DeleteLine(index int) {
	copy(b.lines[index:], b.lines[index+1:])
	b.lines = b.lines[:len(b.lines)-1]
}

func (b *Buffer) DeleteRange(line, start, end int) {
	l := &b.lines[line]

	l.DeleteRange(start, end)
}

// SplitLine splits the specified line at the given column.
//
// The text before the column remains in the original line.
// The text after the column becomes a new line inserted immediately after it.
func (b *Buffer) SplitLine(line, column int) {
	l := &b.lines[line]

	left := l.String()
	runes := []rune(left)

	before := string(runes[:column])
	after := string(runes[column:])

	b.lines[line] = newLine(before)
	b.InsertLine(line + 1)
	b.lines[line+1] = newLine(after)
}

// JoinLines joins the line at index with the following line.
//
// The contents of the following line are appended to the current line.
// The following line is then removed.
func (b *Buffer) JoinLines(line int) {
	current := b.lines[line].String()
	next := b.lines[line+1].String()

	b.lines[line] = newLine(current + next)
	b.DeleteLine(line + 1)
}
