package editor

import (
	"testing"

	"github.com/i-walk-away/daun/internal/buffer"
)

func TestEditorUndoInsert(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	e.InsertText("hello")

	if got := e.Line(0); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}

	e.Undo()

	if got := e.Line(0); got != "" {
		t.Fatalf("expected empty line, got %q", got)
	}

	if got := e.Cursor(); got != (Position{}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{}, got)
	}
}

func TestEditorRedoInsert(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	e.InsertText("hello")
	e.Undo()
	e.Redo()

	if got := e.Line(0); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}

	if got := e.Cursor(); got != (Position{Line: 0, Column: 5}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 5}, got)
	}
}

func TestEditorUndoDelete(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	e.InsertText("hello")
	e.Backspace()

	if got := e.Line(0); got != "hell" {
		t.Fatalf("expected %q, got %q", "hell", got)
	}

	e.Undo()

	if got := e.Line(0); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestEditorUndoMultilineInsert(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	e.InsertText("hello\nworld")

	if got := e.LineCount(); got != 2 {
		t.Fatalf("expected 2 lines, got %d", got)
	}

	e.Undo()

	if got := e.LineCount(); got != 1 {
		t.Fatalf("expected 1 line, got %d", got)
	}

	if got := e.Line(0); got != "" {
		t.Fatalf("expected empty line, got %q", got)
	}
}

func TestEditorUndoDeleteSelection(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	e.InsertText("hello world")
	e.SelectAll()
	e.Backspace()

	if got := e.Line(0); got != "" {
		t.Fatalf("expected empty line, got %q", got)
	}

	e.Undo()

	if got := e.Line(0); got != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", got)
	}
}

func TestEditorUndoRedoMultilineSelection(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	e.InsertText("hello\nbeautiful\nworld")
	e.SelectAll()
	e.Backspace()

	if got := e.LineCount(); got != 1 {
		t.Fatalf("expected 1 line, got %d", got)
	}

	e.Undo()

	if got := e.LineCount(); got != 3 {
		t.Fatalf("expected 3 lines, got %d", got)
	}

	if got := e.Line(1); got != "beautiful" {
		t.Fatalf("expected %q, got %q", "beautiful", got)
	}

	e.Redo()

	if got := e.LineCount(); got != 1 {
		t.Fatalf("expected 1 line after redo, got %d", got)
	}
}

func TestEditorNewEditClearsRedo(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	e.InsertText("hello")
	e.Undo()

	e.InsertText("world")
	e.Redo()

	if got := e.Line(0); got != "world" {
		t.Fatalf("expected %q, got %q", "world", got)
	}
}
