package editor

import (
	"testing"

	"github.com/i-walk-away/daun/internal/buffer"
)

func insertText(b *buffer.Buffer, line int, text string) {
	for _, r := range text {
		b.Insert(line, b.LineLength(line), r)
	}
}

func TestEditorMoveLeft(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")

	e := NewEditor(b)

	e.MoveRight()
	e.MoveRight()
	e.MoveRight()
	e.MoveLeft()

	if got := e.Cursor(); got != (Position{Line: 0, Column: 2}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 2}, got)
	}
}

func TestEditorMoveLeftAtLineStart(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")
	b.InsertLine(1)
	insertText(b, 1, "world")

	e := NewEditor(b)
	e.MoveDown()
	e.MoveLeft()

	if got := e.Cursor(); got != (Position{Line: 0, Column: 5}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 5}, got)
	}
}

func TestEditorMoveLeftAtBufferStart(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")

	e := NewEditor(b)
	e.MoveLeft()

	if got := e.Cursor(); got != (Position{Line: 0, Column: 0}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 0}, got)
	}
}

func TestEditorMoveRight(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")

	e := NewEditor(b)
	e.MoveRight()

	if got := e.Cursor(); got != (Position{Line: 0, Column: 1}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 1}, got)
	}
}

func TestEditorMoveRightAtLineEnd(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")
	b.InsertLine(1)
	insertText(b, 1, "world")

	e := NewEditor(b)

	for range 5 {
		e.MoveRight()
	}

	e.MoveRight()

	if got := e.Cursor(); got != (Position{Line: 1, Column: 0}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 1, Column: 0}, got)
	}
}

func TestEditorMoveRightAtBufferEnd(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")

	e := NewEditor(b)

	for range 6 {
		e.MoveRight()
	}

	if got := e.Cursor(); got != (Position{Line: 0, Column: 5}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 5}, got)
	}
}

func TestEditorMoveUp(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")
	b.InsertLine(1)
	insertText(b, 1, "world")

	e := NewEditor(b)
	e.MoveDown()

	e.MoveRight()
	e.MoveRight()
	e.MoveRight()
	e.MoveUp()

	if got := e.Cursor(); got != (Position{Line: 0, Column: 3}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 3}, got)
	}
}

func TestEditorMoveUpToShorterLine(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hi")
	b.InsertLine(1)
	insertText(b, 1, "hello")

	e := NewEditor(b)
	e.MoveDown()

	for range 5 {
		e.MoveRight()
	}

	e.MoveUp()

	if got := e.Cursor(); got != (Position{Line: 0, Column: 2}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 2}, got)
	}
}

func TestEditorMoveUpAtFirstLine(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")

	e := NewEditor(b)

	for range 3 {
		e.MoveRight()
	}

	e.MoveUp()

	if got := e.Cursor(); got != (Position{Line: 0, Column: 0}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 0}, got)
	}
}

func TestEditorMoveDown(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")
	b.InsertLine(1)
	insertText(b, 1, "world")

	e := NewEditor(b)

	e.MoveRight()
	e.MoveRight()
	e.MoveRight()
	e.MoveDown()

	if got := e.Cursor(); got != (Position{Line: 1, Column: 3}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 1, Column: 3}, got)
	}
}

func TestEditorMoveDownToShorterLine(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")
	b.InsertLine(1)
	insertText(b, 1, "hi")

	e := NewEditor(b)

	for range 5 {
		e.MoveRight()
	}

	e.MoveDown()

	if got := e.Cursor(); got != (Position{Line: 1, Column: 2}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 1, Column: 2}, got)
	}
}

func TestEditorMoveDownAtLastLine(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")

	e := NewEditor(b)

	e.MoveRight()
	e.MoveRight()
	e.MoveDown()

	if got := e.Cursor(); got != (Position{Line: 0, Column: 5}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 5}, got)
	}
}

// Ctrl+Left

func TestEditorMoveWordLeft(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello beautiful world")

	e := NewEditor(b)

	// Start at the end of the buffer.
	for range len([]rune("hello beautiful world")) {
		e.MoveRight()
	}

	// world -> beginning of world.
	e.MoveWordLeft()

	if got := e.Cursor().Column; got != 16 {
		t.Fatalf("expected column 16, got %d", got)
	}

	// beautiful -> beginning of beautiful.
	e.MoveWordLeft()

	if got := e.Cursor().Column; got != 6 {
		t.Fatalf("expected column 6, got %d", got)
	}

	// hello -> beginning of hello.
	e.MoveWordLeft()

	if got := e.Cursor().Column; got != 0 {
		t.Fatalf("expected column 0, got %d", got)
	}
}

func TestEditorMoveWordLeftAcrossWhitespace(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello   beautiful world")

	e := NewEditor(b)

	// Move to the end of the line.
	for range len([]rune("hello   beautiful world")) {
		e.MoveRight()
	}

	// Ctrl+Left from the end of "world":
	// move to the beginning of "world".
	e.MoveWordLeft()

	if got := e.Cursor().Column; got != 18 {
		t.Fatalf("expected cursor column 18, got %d", got)
	}

	// Ctrl+Left from the end of "world":
	// skip whitespace and move to the beginning of "beautiful".
	e.MoveWordLeft()

	if got := e.Cursor().Column; got != 8 {
		t.Fatalf("expected cursor column 8, got %d", got)
	}

	// Ctrl+Left again:
	// move to the beginning of "hello".
	e.MoveWordLeft()

	if got := e.Cursor().Column; got != 0 {
		t.Fatalf("expected cursor column 0, got %d", got)
	}
}

func TestEditorMoveWordLeftInsideWord(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello beautiful")

	e := NewEditor(b)

	for range 10 {
		e.MoveRight()
	}

	e.MoveWordLeft()

	if got := e.Cursor().Column; got != 6 {
		t.Fatalf("expected column 6, got %d", got)
	}
}

func TestEditorMoveWordLeftFromWhitespace(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello   beautiful")

	e := NewEditor(b)

	for range 7 {
		e.MoveRight()
	}

	e.MoveWordLeft()

	if got := e.Cursor().Column; got != 0 {
		t.Fatalf("expected column 0, got %d", got)
	}
}

// Ctrl+Right

func TestEditorMoveWordRight(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello beautiful world")

	e := NewEditor(b)

	// Start of hello -> end of hello.
	e.MoveWordRight()

	if got := e.Cursor().Column; got != 5 {
		t.Fatalf("expected column 5, got %d", got)
	}

	// Skip whitespace, then move to the end of beautiful.
	e.MoveWordRight()

	if got := e.Cursor().Column; got != 15 {
		t.Fatalf("expected column 15, got %d", got)
	}

	// Skip whitespace, then move to the end of world.
	e.MoveWordRight()

	if got := e.Cursor().Column; got != 21 {
		t.Fatalf("expected column 21, got %d", got)
	}
}

func TestEditorMoveWordRightAcrossWhitespace(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello   beautiful world")

	e := NewEditor(b)

	e.MoveWordRight()

	if got := e.Cursor().Column; got != 5 {
		t.Fatalf("expected column 5, got %d", got)
	}

	e.MoveWordRight()

	if got := e.Cursor().Column; got != 17 {
		t.Fatalf("expected column 17, got %d", got)
	}

	e.MoveWordRight()

	if got := e.Cursor().Column; got != 23 {
		t.Fatalf("expected column 23, got %d", got)
	}
}

func TestEditorMoveWordRightInsideWord(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello beautiful")

	e := NewEditor(b)

	for range 8 {
		e.MoveRight()
	}

	e.MoveWordRight()

	if got := e.Cursor().Column; got != 15 {
		t.Fatalf("expected column 15, got %d", got)
	}
}

func TestEditorMoveWordRightFromWhitespace(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello   beautiful")

	e := NewEditor(b)

	for range 6 {
		e.MoveRight()
	}

	e.MoveWordRight()

	if got := e.Cursor().Column; got != 17 {
		t.Fatalf("expected column 17, got %d", got)
	}
}

func TestEditorMoveWordRightAtWordEnd(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello beautiful")

	e := NewEditor(b)

	for range 5 {
		e.MoveRight()
	}

	e.MoveWordRight()

	if got := e.Cursor().Column; got != 15 {
		t.Fatalf("expected column 15, got %d", got)
	}
}

// Cross-line movement

func TestEditorMoveWordLeftAcrossLines(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")
	b.InsertLine(1)
	insertText(b, 1, "world")

	e := NewEditor(b)
	e.MoveDown()
	e.MoveWordLeft()

	if got := e.Cursor(); got != (Position{Line: 0, Column: 5}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 5}, got)
	}
}

func TestEditorMoveWordRightAcrossLines(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")
	b.InsertLine(1)
	insertText(b, 1, "world")

	e := NewEditor(b)

	for range 5 {
		e.MoveRight()
	}

	e.MoveWordRight()

	if got := e.Cursor(); got != (Position{Line: 1, Column: 0}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 1, Column: 0}, got)
	}
}
