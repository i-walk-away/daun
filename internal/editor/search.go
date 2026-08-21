package editor

// Match represents a text match in the document.
//
// Start is inclusive and End is exclusive.
type Match struct {
	Start Position
	End   Position
}

// SearchResult describes the current text-search state.
type SearchResult struct {
	Match   Match
	Current int
	Total   int
	Wrapped bool
}

// Search searches for query starting at from.
//
// The editor state is not modified.
func (e *Editor) Search(
	query string,
	from Position,
) (SearchResult, bool) {
	start := e.clampPosition(from)

	result := e.buffer.Search(
		start.Line,
		start.Column,
		query,
	)

	if result.Total == 0 {
		return SearchResult{}, false
	}

	return SearchResult{
		Match: Match{
			Start: Position{
				Line:   result.StartLine,
				Column: result.StartColumn,
			},
			End: Position{
				Line:   result.EndLine,
				Column: result.EndColumn,
			},
		},
		Current: result.Current,
		Total:   result.Total,
		Wrapped: result.Wrapped,
	}, true
}

// Find searches forward for query starting at from.
func (e *Editor) Find(
	query string,
	from Position,
) (Position, bool) {
	start := e.clampPosition(from)

	line, column, ok := e.buffer.Find(
		start.Line,
		start.Column,
		query,
	)

	if !ok {
		return Position{}, false
	}

	return Position{
		Line:   line,
		Column: column,
	}, true
}

// FindPrevious searches backwards for query starting before from.
func (e *Editor) FindPrevious(
	query string,
	from Position,
) (Position, bool) {
	start := e.clampPosition(from)

	line, column, ok := e.buffer.FindPrevious(
		start.Line,
		start.Column,
		query,
	)

	if !ok {
		return Position{}, false
	}

	return Position{
		Line:   line,
		Column: column,
	}, true
}

// AdvancePosition returns the position immediately after text inserted
// at start.
func AdvancePosition(
	start Position,
	text string,
) Position {
	return advancePosition(start, text)
}
