package buffer

import (
	"strings"
	"testing"
)

func TestBufferFind(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		startLine   int
		startColumn int
		query       string
		wantLine    int
		wantColumn  int
		wantFound   bool
	}{
		{
			name:       "beginning",
			text:       "hello world",
			query:      "hello",
			wantLine:   0,
			wantColumn: 0,
			wantFound:  true,
		},
		{
			name:        "after cursor",
			text:        "hello world hello",
			startColumn: 6,
			query:       "hello",
			wantLine:    0,
			wantColumn:  12,
			wantFound:   true,
		},
		{
			name:       "unicode",
			text:       "Привет мир",
			query:      "мир",
			wantLine:   0,
			wantColumn: 7,
			wantFound:  true,
		},
		{
			name:        "multiline",
			text:        "hello\nbeautiful\nworld",
			startLine:   0,
			startColumn: 0,
			query:       "beautiful",
			wantLine:    1,
			wantColumn:  0,
			wantFound:   true,
		},
		{
			name:       "across newline",
			text:       "hello\nworld",
			query:      "lo\nwo",
			wantLine:   0,
			wantColumn: 3,
			wantFound:  true,
		},
		{
			name:      "missing",
			text:      "hello world",
			query:     "missing",
			wantFound: false,
		},
		{
			name:      "empty query",
			text:      "hello world",
			query:     "",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBufferFromText(tt.text)

			line, column, ok := b.Find(
				tt.startLine,
				tt.startColumn,
				tt.query,
			)

			if ok != tt.wantFound {
				t.Fatalf(
					"expected found=%v, got %v",
					tt.wantFound,
					ok,
				)
			}

			if !tt.wantFound {
				return
			}

			if line != tt.wantLine || column != tt.wantColumn {
				t.Fatalf(
					"expected position (%d, %d), got (%d, %d)",
					tt.wantLine,
					tt.wantColumn,
					line,
					column,
				)
			}
		})
	}
}

func TestBufferFindAcrossPieces(t *testing.T) {
	b := NewBuffer()

	b.InsertText(0, 0, "ind")
	b.InsertText(0, 3, "epen")
	b.InsertText(0, 7, "dence")

	line, column, ok := b.Find(
		0,
		0,
		"independence",
	)

	if !ok {
		t.Fatal("expected match")
	}

	if line != 0 || column != 0 {
		t.Fatalf(
			"expected position (0, 0), got (%d, %d)",
			line,
			column,
		)
	}
}

func TestBufferFindPrevious(t *testing.T) {
	b := NewBufferFromText("one two one two one")

	line, column, ok := b.FindPrevious(
		0,
		14,
		"one",
	)

	if !ok {
		t.Fatal("expected match")
	}

	if line != 0 || column != 8 {
		t.Fatalf(
			"expected position (0, 8), got (%d, %d)",
			line,
			column,
		)
	}
}

func TestBufferSearch(t *testing.T) {
	b := NewBufferFromText("one two one two one")

	result := b.Search(0, 8, "one")

	if result.Total != 3 {
		t.Fatalf(
			"expected total 3, got %d",
			result.Total,
		)
	}

	if result.Current != 2 {
		t.Fatalf(
			"expected current 2, got %d",
			result.Current,
		)
	}

	if result.StartLine != 0 ||
		result.StartColumn != 8 {

		t.Fatalf(
			"expected start (0, 8), got (%d, %d)",
			result.StartLine,
			result.StartColumn,
		)
	}

	if result.EndLine != 0 ||
		result.EndColumn != 11 {

		t.Fatalf(
			"expected end (0, 11), got (%d, %d)",
			result.EndLine,
			result.EndColumn,
		)
	}

	if result.Wrapped {
		t.Fatal("did not expect wrapped search")
	}
}

func TestBufferSearchWraps(t *testing.T) {
	b := NewBufferFromText("one two one")

	result := b.Search(0, 11, "one")

	if result.Total != 2 {
		t.Fatalf(
			"expected total 2, got %d",
			result.Total,
		)
	}

	if result.Current != 1 {
		t.Fatalf(
			"expected current 1, got %d",
			result.Current,
		)
	}

	if !result.Wrapped {
		t.Fatal("expected wrapped search")
	}

	if result.StartLine != 0 ||
		result.StartColumn != 0 {

		t.Fatalf(
			"expected wrapped start (0, 0), got (%d, %d)",
			result.StartLine,
			result.StartColumn,
		)
	}
}

func TestBufferSearchAcrossManyPieces(t *testing.T) {
	b := NewBuffer()

	parts := []string{
		"alpha",
		" beta",
		" gamma",
		" delta",
	}

	column := 0

	for _, part := range parts {
		b.InsertText(0, column, part)
		column += len([]rune(part))
	}

	result := b.Search(0, 0, "beta gamma")

	if result.Total != 1 {
		t.Fatalf(
			"expected one match, got %d",
			result.Total,
		)
	}

	if result.StartLine != 0 ||
		result.StartColumn != 6 {

		t.Fatalf(
			"expected start (0, 6), got (%d, %d)",
			result.StartLine,
			result.StartColumn,
		)
	}
}

func TestBufferSearchLargeDocument(t *testing.T) {
	const lineCount = 100_000

	text := strings.Repeat(
		"hello world\n",
		lineCount,
	)

	b := NewBufferFromText(text)

	lastLine := lineCount - 1

	b.InsertText(
		lastLine,
		b.LineLength(lastLine),
		"TARGET",
	)

	result := b.Search(0, 0, "TARGET")

	if result.Total != 1 {
		t.Fatalf(
			"expected one match, got %d",
			result.Total,
		)
	}

	if result.StartLine != lastLine {
		t.Fatalf(
			"expected line %d, got %d",
			lastLine,
			result.StartLine,
		)
	}

	if result.StartColumn != 11 {
		t.Fatalf(
			"expected column 11, got %d",
			result.StartColumn,
		)
	}
}
