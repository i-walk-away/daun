package buffer

const initialGapSize = 16

// gapBuffer stores a sequence of runes with an unused region between
// gapStart and gapEnd.
//
// The gap represents the current editing position. Inserting or deleting
// characters at the gap does not require moving the rest of the text.
type gapBuffer struct {
	data     []rune
	gapStart int
	gapEnd   int
}

// newGapBuffer creates a gap buffer containing text.
//
// The gap is initially placed at the end of the text.
func newGapBuffer(text string) gapBuffer {
	runes := []rune(text)
	data := make([]rune, len(runes)+initialGapSize)

	copy(data, runes)

	return gapBuffer{
		data:     data,
		gapStart: len(runes),
		gapEnd:   len(data),
	}
}

// Len returns the number of runes stored in the gap buffer.
func (g *gapBuffer) Len() int {
	return len(g.data) - (g.gapEnd - g.gapStart)
}

// Insert inserts a rune at the current gap position.
func (g *gapBuffer) Insert(r rune) {
	if g.gapStart == g.gapEnd {
		g.grow()
	}

	g.data[g.gapStart] = r
	g.gapStart++
}

// Delete removes the rune immediately before the current gap position.
//
// Delete does nothing when the gap is at the beginning of the buffer.
func (g *gapBuffer) Delete() {
	if g.gapStart == 0 {
		return
	}

	g.gapStart--
}

// DeleteRange removes the runes in the half-open range [start, end).
//
// The start index is inclusive and the end index is exclusive.
// For example, DeleteRange(2, 5) removes runes at positions 2, 3, and 4.
func (g *gapBuffer) DeleteRange(start, end int) {
	if start < 0 || end > g.Len() || start >= end {
		return
	}

	g.MoveGap(start)

	deleteLength := end - start
	g.gapEnd += deleteLength
}

// MoveGap moves the gap to the specified logical position.
//
// Position is measured in runes from the beginning of the logical text.
// Moving the gap does not change the logical contents of the buffer.
func (g *gapBuffer) MoveGap(position int) {
	if position < 0 || position > g.Len() {
		return
	}

	switch {
	case position < g.gapStart:
		g.moveGapLeft(position)

	case position > g.gapStart:
		g.moveGapRight(position)
	}
}

// String returns the logical contents of the gap buffer.
//
// The unused gap is not included in the returned string.
func (g *gapBuffer) String() string {
	result := make([]rune, 0, g.Len())

	result = append(result, g.data[:g.gapStart]...)
	result = append(result, g.data[g.gapEnd:]...)

	return string(result)
}

// grow increases the size of the gap while preserving the logical contents.
func (g *gapBuffer) grow() {
	newData := make([]rune, len(g.data)+initialGapSize)

	copy(newData, g.data[:g.gapStart])

	newGapEnd := g.gapEnd + initialGapSize

	copy(newData[newGapEnd:], g.data[g.gapEnd:])

	g.data = newData
	g.gapEnd = newGapEnd
}

// moveGapLeft moves the gap to a position before its current position.
func (g *gapBuffer) moveGapLeft(position int) {
	distance := g.gapStart - position

	copy(
		g.data[g.gapEnd-distance:g.gapEnd],
		g.data[position:g.gapStart],
	)

	g.gapStart -= distance
	g.gapEnd -= distance
}

// moveGapRight moves the gap to a position after its current position.
func (g *gapBuffer) moveGapRight(position int) {
	distance := position - g.gapStart

	copy(
		g.data[g.gapStart:g.gapStart+distance],
		g.data[g.gapEnd:g.gapEnd+distance],
	)

	g.gapStart += distance
	g.gapEnd += distance
}
