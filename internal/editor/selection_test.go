package editor

import (
	"testing"

	"github.com/i-walk-away/daun/internal/buffer"
)

func TestEditorSelectAll(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")

	b.InsertLine(1)
	insertText(b, 1, "beautiful")

	b.InsertLine(2)
	insertText(b, 2, "world")

	e := NewEditor(b)
	e.SelectAll()

	start, end, ok := e.Selection()
	if !ok {
		t.Fatal("expected active selection")
	}

	if start != (Position{Line: 0, Column: 0}) {
		t.Fatalf("expected start %+v, got %+v", Position{}, start)
	}

	expectedEnd := Position{
		Line:   2,
		Column: 5,
	}

	if end != expectedEnd {
		t.Fatalf("expected end %+v, got %+v", expectedEnd, end)
	}

	text, ok := e.SelectedText()
	if !ok {
		t.Fatal("expected selected text")
	}

	if text != "hello\nbeautiful\nworld" {
		t.Fatalf("expected selected text %q, got %q", "hello\nbeautiful\nworld", text)
	}
}

func TestEditorSelectAllEmptyBuffer(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	e.SelectAll()

	_, _, ok := e.Selection()
	if ok {
		t.Fatal("expected no selection for an empty buffer")
	}
}

func TestEditorSelectedTextSingleLine(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello world")

	e := NewEditor(b)

	e.Move(
		DirectionRight,
		MoveByCharacter,
		true,
	)
	e.Move(
		DirectionRight,
		MoveByCharacter,
		true,
	)
	e.Move(
		DirectionRight,
		MoveByCharacter,
		true,
	)

	text, ok := e.SelectedText()
	if !ok {
		t.Fatal("expected active selection")
	}

	if text != "hel" {
		t.Fatalf("expected %q, got %q", "hel", text)
	}
}

func TestEditorSelectedTextMultipleLines(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")

	b.InsertLine(1)
	insertText(b, 1, "beautiful")

	b.InsertLine(2)
	insertText(b, 2, "world")

	e := NewEditor(b)

	for range 2 {
		e.Move(
			DirectionDown,
			MoveByCharacter,
			true,
		)
	}

	text, ok := e.SelectedText()
	if !ok {
		t.Fatal("expected active selection")
	}

	if text != "hello\nbeautiful\n" {
		t.Fatalf("expected %q, got %q", "hello\nbeautiful\n", text)
	}
}

func TestEditorSelectionCanMoveBackwards(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")

	e := NewEditor(b)

	e.Move(
		DirectionRight,
		MoveByCharacter,
		true,
	)
	e.Move(
		DirectionRight,
		MoveByCharacter,
		true,
	)
	e.Move(
		DirectionRight,
		MoveByCharacter,
		true,
	)

	e.Move(
		DirectionLeft,
		MoveByCharacter,
		true,
	)

	start, end, ok := e.Selection()
	if !ok {
		t.Fatal("expected active selection")
	}

	if start != (Position{Line: 0, Column: 0}) {
		t.Fatalf("expected start %+v, got %+v", Position{}, start)
	}

	if end != (Position{Line: 0, Column: 2}) {
		t.Fatalf(
			"expected end %+v, got %+v",
			Position{Line: 0, Column: 2},
			end,
		)
	}
}

func TestEditorClearSelection(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")

	e := NewEditor(b)

	e.SelectAll()
	e.ClearSelection()

	_, _, ok := e.Selection()
	if ok {
		t.Fatal("expected selection to be cleared")
	}
}
