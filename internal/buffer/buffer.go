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

// DeleteRange removes the range [start, end) from the specified line.
func (b *Buffer) DeleteRange(line, start, end int) {
	l := &b.lines[line]

	l.DeleteRange(start, end)
}

// DeleteRangeLines removes the text in the range [start, end).
//
// The range may span multiple lines. Text before start on the first line
// and text after end on the last line are preserved and joined together.
func (b *Buffer) DeleteRangeLines(
	startLine, startColumn,
	endLine, endColumn int,
) {
	if startLine == endLine {
		b.DeleteRange(startLine, startColumn, endColumn)
		return
	}

	first := []rune(b.Line(startLine))
	last := []rune(b.Line(endLine))

	startColumn = min(max(startColumn, 0), len(first))
	endColumn = min(max(endColumn, 0), len(last))

	merged := make([]rune, 0, startColumn+len(last)-endColumn)
	merged = append(merged, first[:startColumn]...)
	merged = append(merged, last[endColumn:]...)

	b.lines[startLine] = newLine(string(merged))

	removedLines := endLine - startLine

	copy(
		b.lines[startLine+1:],
		b.lines[endLine+1:],
	)

	b.lines = b.lines[:len(b.lines)-removedLines]
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
