package tui

import (
	"strings"
	"testing"

	"github.com/i-walk-away/daun/internal/editor"
)

func TestRenderLineSearchPreservesText(t *testing.T) {
	text := "hello hello hello"

	result := renderLineSearch(
		text,
		0,
		"hello",
		nil,
		0,
		100,
	)

	plain := stripANSI(result)

	if plain != text {
		t.Fatalf(
			"expected %q, got %q",
			text,
			plain,
		)
	}
}

func TestRenderLineSearchPreservesUnicode(t *testing.T) {
	text := "Привет hello мир hello"

	result := renderLineSearch(
		text,
		0,
		"hello",
		nil,
		0,
		100,
	)

	plain := stripANSI(result)

	if plain != text {
		t.Fatalf(
			"expected %q, got %q",
			text,
			plain,
		)
	}
}

func TestRenderLineSearchRespectsHorizontalViewport(t *testing.T) {
	result := renderLineSearch(
		"0123456789 hello world",
		0,
		"hello",
		nil,
		10,
		5,
	)

	plain := stripANSI(result)

	if plain != " hell" {
		t.Fatalf(
			"expected %q, got %q",
			" hell",
			plain,
		)
	}
}

func TestRenderLineSearchClipsToContentWidth(t *testing.T) {
	result := renderLineSearch(
		"0123456789hello world",
		0,
		"hello",
		nil,
		10,
		5,
	)

	plain := stripANSI(result)

	if plain != "hello" {
		t.Fatalf(
			"expected %q, got %q",
			"hello",
			plain,
		)
	}
}

func TestRenderLineSearchHighlightsCurrentMatchWithoutChangingText(
	t *testing.T,
) {
	current := &editor.Match{
		Start: editor.Position{
			Line:   0,
			Column: 6,
		},
		End: editor.Position{
			Line:   0,
			Column: 11,
		},
	}

	text := "hello hello"

	result := renderLineSearch(
		text,
		0,
		"hello",
		current,
		0,
		100,
	)

	plain := stripANSI(result)

	if plain != text {
		t.Fatalf(
			"expected %q, got %q",
			text,
			plain,
		)
	}
}

func TestRenderLineSearchDoesNotRenderMultilineQueryOnSingleLine(
	t *testing.T,
) {
	text := "hello world"

	result := renderLineSearch(
		text,
		0,
		"hello\nworld",
		nil,
		0,
		100,
	)

	plain := stripANSI(result)

	if plain != text {
		t.Fatalf(
			"expected %q, got %q",
			text,
			plain,
		)
	}
}

func TestRenderLineSearchEmptyQuery(t *testing.T) {
	text := "hello world"

	result := renderLineSearch(
		text,
		0,
		"",
		nil,
		0,
		100,
	)

	if result != text {
		t.Fatalf(
			"expected %q, got %q",
			text,
			result,
		)
	}
}

func stripANSI(text string) string {
	var builder strings.Builder

	inEscape := false

	for _, r := range text {
		switch {
		case r == '\x1b':
			inEscape = true

		case inEscape && r == 'm':
			inEscape = false

		case !inEscape:
			builder.WriteRune(r)
		}
	}

	return builder.String()
}
