package editor

import (
	"testing"

	"github.com/i-walk-away/daun/internal/buffer"
)

func TestEditorSetCursor(t *testing.T) {
	b := buffer.NewBufferFromText("hello\nworld")
	e := NewEditor(b)

	e.SetCursor(Position{
		Line:   1,
		Column: 3,
	})

	if got := e.Cursor(); got != (Position{Line: 1, Column: 3}) {
		t.Fatalf(
			"expected cursor %+v, got %+v",
			Position{Line: 1, Column: 3},
			got,
		)
	}
}

func TestEditorSetCursorClampsPosition(t *testing.T) {
	b := buffer.NewBufferFromText("hello")
	e := NewEditor(b)

	e.SetCursor(Position{
		Line:   100,
		Column: 100,
	})

	if got := e.Cursor(); got != (Position{Line: 0, Column: 5}) {
		t.Fatalf(
			"expected cursor %+v, got %+v",
			Position{Line: 0, Column: 5},
			got,
		)
	}
}

func TestEditorFind(t *testing.T) {
	b := buffer.NewBufferFromText(
		"hello\nbeautiful\nworld",
	)

	e := NewEditor(b)

	position, ok := e.Find(
		"beautiful",
		Position{},
	)

	if !ok {
		t.Fatal("expected match")
	}

	if position != (Position{Line: 1, Column: 0}) {
		t.Fatalf(
			"expected position %+v, got %+v",
			Position{Line: 1, Column: 0},
			position,
		)
	}

	if e.Cursor() != (Position{}) {
		t.Fatal("Find must not move the cursor")
	}
}
