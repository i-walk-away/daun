package buffer

import (
	"strings"
	"unicode/utf8"
)

// Buffer stores editor text in a piece table backed by a balanced B+ tree.
// The tree indexes byte, rune, and newline counts for logarithmic document
// navigation while keeping inserted and original text in append-only storage.
type Buffer struct {
	sources *sources
	tree    *pieceTree
}

// NewBuffer creates an empty buffer containing one empty line.
func NewBuffer() *Buffer {
	s := &sources{}

	return &Buffer{
		sources: s,
		tree:    newPieceTree(s, nil),
	}
}

// Close releases resources owned by the buffer.
//
// For buffers created with NewBufferFromFile, this releases the underlying
// memory mapping. Buffers created with NewBuffer do not own external resources.
func (b *Buffer) Close() error {
	if b == nil || b.sources == nil {
		return nil
	}

	return b.sources.close()
}

// NewBufferFromText creates a buffer initialized with text.
func NewBufferFromText(text string) *Buffer {
	text = normalizeText(text)

	s := &sources{
		original: []byte(text),
	}

	pieces := makePieces(sourceOriginal, s.original)

	return &Buffer{
		sources: s,
		tree:    newPieceTree(s, pieces),
	}
}

// Line returns the text of the line at index.
func (b *Buffer) Line(index int) string {
	start, end := b.lineByteRange(index)

	var builder strings.Builder
	builder.Grow(end - start)

	b.tree.readRange(
		start,
		end,
		func(data []byte) {
			_, _ = builder.Write(data)
		},
	)

	return builder.String()
}

// LineLength returns the number of runes in the line at index.
func (b *Buffer) LineLength(index int) int {
	start, end := b.lineByteRange(index)

	startMetrics := b.tree.prefixMetrics(start)
	endMetrics := b.tree.prefixMetrics(end)

	return endMetrics.runes - startMetrics.runes
}

// LineCount returns the number of lines in the buffer.
func (b *Buffer) LineCount() int {
	return b.tree.totalNewlines() + 1
}

// Insert inserts a rune at the specified line and column.
func (b *Buffer) Insert(line, column int, r rune) {
	b.InsertText(line, column, string(r))
}

// InsertText inserts text at the specified line and column.
func (b *Buffer) InsertText(line, column int, text string) {
	text = normalizeText(text)

	if text == "" {
		return
	}

	offset := b.positionToOffset(line, column)

	addStart := len(b.sources.add)
	b.sources.add = append(b.sources.add, []byte(text)...)

	for _, p := range makePieces(sourceAdd, b.sources.add[addStart:]) {
		p.start += addStart
		p.end += addStart

		b.tree.insertPiece(offset, p)
		offset += p.metrics.bytes
	}
}

// Delete removes the rune immediately before the specified position.
func (b *Buffer) Delete(line, column int) {
	if line < 0 || line >= b.LineCount() || column <= 0 {
		return
	}

	start := b.positionToOffset(line, column-1)
	end := b.positionToOffset(line, column)

	b.tree.deleteRange(start, end)
}

// InsertLine inserts an empty line at the specified index.
func (b *Buffer) InsertLine(index int) {
	lineCount := b.LineCount()
	index = min(max(index, 0), lineCount)

	if index == lineCount {
		last := lineCount - 1
		b.InsertText(last, b.LineLength(last), "\n")
		return
	}

	b.InsertText(index, 0, "\n")
}

// DeleteLine removes the line at the specified index.
//
// The buffer always retains at least one line.
func (b *Buffer) DeleteLine(index int) {
	lineCount := b.LineCount()

	if lineCount <= 1 || index < 0 || index >= lineCount {
		return
	}

	switch {
	case index < lineCount-1:
		start := b.positionToOffset(index, 0)
		end := b.positionToOffset(index+1, 0)

		b.tree.deleteRange(start, end)

	default:
		start := b.positionToOffset(index-1, b.LineLength(index-1))
		end := b.positionToOffset(index, b.LineLength(index))

		b.tree.deleteRange(start, end)
	}
}

// DeleteRange removes the range [start, end) from a single line.
func (b *Buffer) DeleteRange(line, start, end int) {
	if line < 0 || line >= b.LineCount() || start >= end {
		return
	}

	startOffset := b.positionToOffset(line, start)
	endOffset := b.positionToOffset(line, end)

	b.tree.deleteRange(startOffset, endOffset)
}

// DeleteRangeLines removes the range [start, end), which may span lines.
func (b *Buffer) DeleteRangeLines(
	startLine, startColumn,
	endLine, endColumn int,
) {
	start := b.positionToOffset(startLine, startColumn)
	end := b.positionToOffset(endLine, endColumn)

	b.tree.deleteRange(start, end)
}

// TextRange returns the text in the range [start, end).
func (b *Buffer) TextRange(
	startLine, startColumn,
	endLine, endColumn int,
) string {
	start := b.positionToOffset(startLine, startColumn)
	end := b.positionToOffset(endLine, endColumn)

	if start >= end {
		return ""
	}

	var builder strings.Builder
	builder.Grow(end - start)

	b.tree.readRange(
		start,
		end,
		func(data []byte) {
			_, _ = builder.Write(data)
		},
	)

	return builder.String()
}

// ReplaceRange replaces the range [start, end) with replacement.
func (b *Buffer) ReplaceRange(
	startLine, startColumn,
	endLine, endColumn int,
	replacement string,
) {
	start := b.positionToOffset(startLine, startColumn)
	end := b.positionToOffset(endLine, endColumn)

	replacement = normalizeText(replacement)

	b.tree.deleteRange(start, end)

	if replacement == "" {
		return
	}

	addStart := len(b.sources.add)
	b.sources.add = append(b.sources.add, []byte(replacement)...)

	for _, p := range makePieces(sourceAdd, b.sources.add[addStart:]) {
		p.start += addStart
		p.end += addStart

		b.tree.insertPiece(start, p)
		start += p.metrics.bytes
	}
}

// SplitLine splits a line at column.
func (b *Buffer) SplitLine(line, column int) {
	b.ReplaceRange(line, column, line, column, "\n")
}

// JoinLines joins the line at index with the following line.
func (b *Buffer) JoinLines(line int) {
	if line < 0 || line >= b.LineCount()-1 {
		return
	}

	column := b.LineLength(line)
	b.DeleteRangeLines(line, column, line+1, 0)
}

func (b *Buffer) lineByteRange(index int) (int, int) {
	lineCount := b.LineCount()
	index = min(max(index, 0), lineCount-1)

	start := b.tree.offsetByNewlines(index)

	var end int

	if index == lineCount-1 {
		end = b.tree.totalBytes()
	} else {
		end = b.tree.offsetByNewlines(index + 1)

		if end > start {
			end--
		}
	}

	return start, max(end, start)
}

func (b *Buffer) positionToOffset(line, column int) int {
	lineCount := b.LineCount()
	line = min(max(line, 0), lineCount-1)

	start, end := b.lineByteRange(line)

	startMetrics := b.tree.prefixMetrics(start)
	endMetrics := b.tree.prefixMetrics(end)

	lineRunes := endMetrics.runes - startMetrics.runes
	column = min(max(column, 0), lineRunes)

	targetRunes := startMetrics.runes + column

	return min(
		b.tree.offsetByRunes(targetRunes),
		end,
	)
}

func makePieces(source sourceKind, data []byte) []piece {
	if len(data) == 0 {
		return nil
	}

	pieces := make(
		[]piece,
		0,
		(len(data)+maxPieceBytes-1)/maxPieceBytes,
	)

	for start := 0; start < len(data); {
		end := min(start+maxPieceBytes, len(data))

		// Never split a UTF-8 sequence between pieces.
		if end < len(data) {
			for end > start && !utf8.RuneStart(data[end]) {
				end--
			}

			if end == start {
				end = min(start+maxPieceBytes, len(data))

				for end < len(data) && !utf8.RuneStart(data[end]) {
					end++
				}
			}
		}

		pieces = append(
			pieces,
			newPiece(source, start, end, data),
		)

		start = end
	}

	return pieces
}

func normalizeText(text string) string {
	return strings.ReplaceAll(text, "\r\n", "\n")
}
