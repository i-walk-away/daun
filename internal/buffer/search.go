package buffer

import (
	"bytes"
	"strings"
)

// SearchResult describes the selected result of a document search.
//
// Start and End use the same zero-based line and rune-column coordinates as
// the rest of the buffer API. End is exclusive.
//
// Current and Total are one-based when a match exists. Wrapped reports
// whether the search had to continue from the beginning of the document.
type SearchResult struct {
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
	Current     int
	Total       int
	Wrapped     bool
}

// Find searches for query starting at line and column.
//
// The search operates directly over the piece table and never materializes
// the complete document. The returned position is the start of the first
// match at or after the supplied position.
func (b *Buffer) Find(
	line int,
	column int,
	query string,
) (foundLine int, foundColumn int, ok bool) {
	query = NormalizeSearchQuery(query)

	if b == nil || b.tree == nil || query == "" {
		return 0, 0, false
	}

	start := b.positionToOffset(line, column)
	queryBytes := []byte(query)

	offset, found := b.tree.findBytes(start, queryBytes)
	if !found {
		return 0, 0, false
	}

	foundLine, foundColumn = b.offsetToPosition(offset)

	return foundLine, foundColumn, true
}

// FindPrevious searches backwards for the last match whose start is before
// the supplied position.
//
// The returned position is the start of the match.
func (b *Buffer) FindPrevious(
	line int,
	column int,
	query string,
) (foundLine int, foundColumn int, ok bool) {
	query = NormalizeSearchQuery(query)

	if b == nil || b.tree == nil || query == "" {
		return 0, 0, false
	}

	end := b.positionToOffset(line, column)
	queryBytes := []byte(query)

	offset, found := b.tree.findPreviousBytes(end, queryBytes)
	if !found {
		return 0, 0, false
	}

	foundLine, foundColumn = b.offsetToPosition(offset)

	return foundLine, foundColumn, true
}

// Search finds all matches logically, without retaining all match positions,
// and selects the first match at or after the supplied position.
//
// When no match occurs at or after the supplied position, Search wraps to the
// first match in the document and sets Wrapped to true.
//
// Search uses bounded memory proportional to the query size rather than the
// document size.
func (b *Buffer) Search(
	line int,
	column int,
	query string,
) SearchResult {
	query = NormalizeSearchQuery(query)

	if b == nil || b.tree == nil || query == "" {
		return SearchResult{}
	}

	from := b.positionToOffset(line, column)
	queryBytes := []byte(query)

	total := 0
	firstOffset := -1
	selectedOffset := -1
	selectedOrdinal := 0

	b.tree.scanMatches(
		0,
		b.tree.totalBytes(),
		queryBytes,
		func(offset int) bool {
			total++

			if firstOffset < 0 {
				firstOffset = offset
			}

			if selectedOffset < 0 && offset >= from {
				selectedOffset = offset
				selectedOrdinal = total
			}

			return false
		},
	)

	if total == 0 {
		return SearchResult{}
	}

	wrapped := false

	if selectedOffset < 0 {
		selectedOffset = firstOffset
		selectedOrdinal = 1
		wrapped = true
	}

	startLine, startColumn := b.offsetToPosition(selectedOffset)
	endLine, endColumn := b.offsetToPosition(
		selectedOffset + len(queryBytes),
	)

	return SearchResult{
		StartLine:   startLine,
		StartColumn: startColumn,
		EndLine:     endLine,
		EndColumn:   endColumn,
		Current:     selectedOrdinal,
		Total:       total,
		Wrapped:     wrapped,
	}
}

// NormalizeSearchQuery normalizes platform-independent line endings to the
// newline representation used internally by the buffer.
func NormalizeSearchQuery(query string) string {
	return strings.ReplaceAll(query, "\r\n", "\n")
}

func (b *Buffer) offsetToPosition(offset int) (int, int) {
	offset = min(max(offset, 0), b.tree.totalBytes())

	metrics := b.tree.prefixMetrics(offset)

	line := metrics.newlines
	lineStart := b.tree.offsetByNewlines(line)
	lineStartMetrics := b.tree.prefixMetrics(lineStart)

	return line, metrics.runes - lineStartMetrics.runes
}

func (t *pieceTree) findBytes(
	start int,
	query []byte,
) (int, bool) {
	if len(query) == 0 {
		return start, false
	}

	found := -1

	t.scanMatches(
		start,
		t.totalBytes(),
		query,
		func(offset int) bool {
			found = offset
			return true
		},
	)

	return found, found >= 0
}

func (t *pieceTree) findPreviousBytes(
	end int,
	query []byte,
) (int, bool) {
	if len(query) == 0 || end <= 0 {
		return 0, false
	}

	last := -1

	t.scanMatches(
		0,
		end,
		query,
		func(offset int) bool {
			last = offset
			return false
		},
	)

	return last, last >= 0
}

// scanMatches visits every match whose start offset lies in [start, end).
//
// The callback may return true to stop the scan.
//
// Each piece is searched with bytes.Index. A bounded carry buffer preserves
// up to len(query)-1 bytes between adjacent pieces, allowing a match to cross
// a piece boundary without materializing the complete document.
func (t *pieceTree) scanMatches(
	start int,
	end int,
	query []byte,
	visit func(offset int) bool,
) {
	if len(query) == 0 {
		return
	}

	start = min(max(start, 0), t.totalBytes())
	end = min(max(end, start), t.totalBytes())

	if start >= end {
		return
	}

	carrySize := len(query) - 1
	carry := make([]byte, 0, carrySize)

	lastReported := start - 1

	t.scanChunks(
		start,
		end,
		func(chunkStart int, data []byte) bool {
			if len(data) == 0 {
				return false
			}

			window := make([]byte, len(carry)+len(data))

			copy(window, carry)
			copy(window[len(carry):], data)

			baseOffset := chunkStart - len(carry)

			searchAt := 0

			for searchAt < len(window) {
				index := bytes.Index(
					window[searchAt:],
					query,
				)
				if index < 0 {
					break
				}

				index += searchAt

				matchOffset := baseOffset + index

				if matchOffset >= start &&
					matchOffset < end &&
					matchOffset > lastReported {

					lastReported = matchOffset

					if visit(matchOffset) {
						return true
					}
				}

				// Advance by one byte so overlapping matches are preserved.
				searchAt = index + 1
			}

			if carrySize == 0 {
				carry = carry[:0]
				return false
			}

			if len(window) <= carrySize {
				carry = append(carry[:0], window...)
				return false
			}

			carry = append(
				carry[:0],
				window[len(window)-carrySize:]...,
			)

			return false
		},
	)
}

func (t *pieceTree) scanChunks(
	start int,
	end int,
	visit func(offset int, data []byte) bool,
) {
	if start >= end || t.totalBytes() == 0 {
		return
	}

	leaf, index, inner := t.locate(start)
	offset := start

	for leaf != nil && offset < end {
		if index >= len(leaf.pieces) {
			leaf = leaf.next
			index = 0
			inner = 0
			continue
		}

		piece := leaf.pieces[index]

		data := t.sources.bytes(
			piece.source,
			piece.start,
			piece.end,
		)

		if inner > 0 {
			data = data[inner:]
		}

		remaining := end - offset

		if len(data) > remaining {
			data = data[:remaining]
		}

		if visit(offset, data) {
			return
		}

		offset += len(data)
		index++
		inner = 0
	}
}
