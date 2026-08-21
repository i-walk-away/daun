package editor

import (
	"testing"

	"github.com/i-walk-away/daun/internal/buffer"
)

func TestEditorInsert(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	e.Insert('h')
	e.Insert('i')

	if got := e.Line(0); got != "hi" {
		t.Fatalf("expected %q, got %q", "hi", got)
	}

	if got := e.Cursor(); got != (Position{Line: 0, Column: 2}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 2}, got)
	}
}

func TestEditorInsertInMiddle(t *testing.T) {
	b := buffer.NewBuffer()

	for _, r := range "helo" {
		b.Insert(0, b.LineLength(0), r)
	}

	e := NewEditor(b)

	// Move between 'e' and 'l'.
	e.moveRight()
	e.moveRight()

	e.Insert('l')

	if got := e.Line(0); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}

	if got := e.Cursor(); got != (Position{Line: 0, Column: 3}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 3}, got)
	}
}

func TestEditorInsertUnicode(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	e.Insert('П')
	e.Insert('р')
	e.Insert('и')
	e.Insert('в')
	e.Insert('е')
	e.Insert('т')

	if got := e.Line(0); got != "Привет" {
		t.Fatalf("expected %q, got %q", "Привет", got)
	}

	if got := e.Cursor().Column; got != 6 {
		t.Fatalf("expected cursor column 6, got %d", got)
	}
}

func TestEditorBackspace(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	for _, r := range "hello" {
		e.Insert(r)
	}

	e.Backspace()

	if got := e.Line(0); got != "hell" {
		t.Fatalf("expected %q, got %q", "hell", got)
	}

	if got := e.Cursor().Column; got != 4 {
		t.Fatalf("expected cursor column 4, got %d", got)
	}
}

func TestEditorBackspaceAtLineStart(t *testing.T) {
	b := buffer.NewBuffer()

	for _, r := range "hello" {
		b.Insert(0, b.LineLength(0), r)
	}

	b.InsertLine(1)

	for _, r := range "world" {
		b.Insert(1, b.LineLength(1), r)
	}

	e := NewEditor(b)

	// Move cursor to the beginning of the second line.
	e.moveDown()

	e.Backspace()

	if got := e.LineCount(); got != 1 {
		t.Fatalf("expected 1 line, got %d", got)
	}

	if got := e.Line(0); got != "helloworld" {
		t.Fatalf("expected %q, got %q", "helloworld", got)
	}

	if got := e.Cursor(); got != (Position{Line: 0, Column: 5}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 5}, got)
	}
}

func TestEditorBackspaceAtBufferStart(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	e.Backspace()

	if got := e.Line(0); got != "" {
		t.Fatalf("expected empty line, got %q", got)
	}

	if got := e.Cursor(); got != (Position{Line: 0, Column: 0}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 0}, got)
	}
}

func TestEditorEnter(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	for _, r := range "hello world" {
		e.Insert(r)
	}

	e.Cursor()

	// Move cursor between "hello" and " world".
	for i := 0; i < 6; i++ {
		e.moveLeft()
	}

	e.Enter()

	if got := e.LineCount(); got != 2 {
		t.Fatalf("expected 2 lines, got %d", got)
	}

	if got := e.Line(0); got != "hello" {
		t.Fatalf("expected first line %q, got %q", "hello", got)
	}

	if got := e.Line(1); got != " world" {
		t.Fatalf("expected second line %q, got %q", " world", got)
	}

	if got := e.Cursor(); got != (Position{Line: 1, Column: 0}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 1, Column: 0}, got)
	}
}

func TestEditorEnterAtLineStart(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	for _, r := range "hello" {
		e.Insert(r)
	}

	e.moveLeft()
	e.moveLeft()
	e.moveLeft()
	e.moveLeft()
	e.moveLeft()

	e.Enter()

	if got := e.Line(0); got != "" {
		t.Fatalf("expected first line empty, got %q", got)
	}

	if got := e.Line(1); got != "hello" {
		t.Fatalf("expected second line %q, got %q", "hello", got)
	}

	if got := e.Cursor(); got != (Position{Line: 1, Column: 0}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 1, Column: 0}, got)
	}
}

func TestEditorEnterAtLineEnd(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	for _, r := range "hello" {
		e.Insert(r)
	}

	e.Enter()

	if got := e.Line(0); got != "hello" {
		t.Fatalf("expected first line %q, got %q", "hello", got)
	}

	if got := e.Line(1); got != "" {
		t.Fatalf("expected second line empty, got %q", got)
	}

	if got := e.Cursor(); got != (Position{Line: 1, Column: 0}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 1, Column: 0}, got)
	}
}

func TestEditorDeleteWordBackwards(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	for _, r := range "hello beautiful world" {
		e.Insert(r)
	}

	// Cursor is at the end of the line.
	e.DeleteWordBackwards()

	if got := e.Line(0); got != "hello beautiful " {
		t.Fatalf("expected %q, got %q", "hello beautiful ", got)
	}

	if got := e.Cursor().Column; got != len([]rune("hello beautiful ")) {
		t.Fatalf(
			"expected cursor column %d, got %d",
			len([]rune("hello beautiful ")),
			e.Cursor().Column,
		)
	}
}

func TestEditorDeleteWordBackwardsFromWhitespace(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	for _, r := range "hello   beautiful" {
		e.Insert(r)
	}

	// Move cursor to the end of the whitespace before "beautiful".
	for i := 0; i < len([]rune("beautiful")); i++ {
		e.moveLeft()
	}

	e.DeleteWordBackwards()

	if got := e.Line(0); got != "beautiful" {
		t.Fatalf("expected %q, got %q", "beautiful", got)
	}
}

func TestEditorDeleteWordBackwardsAtStartOfLine(t *testing.T) {
	b := buffer.NewBuffer()

	for _, r := range "hello" {
		b.Insert(0, b.LineLength(0), r)
	}

	b.InsertLine(1)

	for _, r := range "world" {
		b.Insert(1, b.LineLength(1), r)
	}

	e := NewEditor(b)
	e.moveDown()

	e.DeleteWordBackwards()

	if got := e.LineCount(); got != 1 {
		t.Fatalf("expected 1 line, got %d", got)
	}

	if got := e.Line(0); got != "helloworld" {
		t.Fatalf("expected %q, got %q", "helloworld", got)
	}

	if got := e.Cursor(); got != (Position{Line: 0, Column: 5}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 5}, got)
	}
}

func TestEditorDeleteWordBackwardsAtBufferStart(t *testing.T) {
	b := buffer.NewBuffer()
	e := NewEditor(b)

	e.DeleteWordBackwards()

	if got := e.Line(0); got != "" {
		t.Fatalf("expected empty line, got %q", got)
	}

	if got := e.Cursor(); got != (Position{Line: 0, Column: 0}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 0}, got)
	}
}

func TestEditorInsertReplacesSelection(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello world")

	e := NewEditor(b)
	e.SelectAll()

	e.Insert('X')

	if got := e.Line(0); got != "X" {
		t.Fatalf("expected %q, got %q", "X", got)
	}

	if got := e.Cursor(); got != (Position{Line: 0, Column: 1}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 1}, got)
	}

	_, _, ok := e.Selection()
	if ok {
		t.Fatal("expected selection to be cleared")
	}
}

func TestEditorBackspaceDeletesSelection(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello world")

	e := NewEditor(b)
	e.SelectAll()

	e.Backspace()

	if got := e.Line(0); got != "" {
		t.Fatalf("expected empty line, got %q", got)
	}

	if got := e.Cursor(); got != (Position{}) {
		t.Fatalf("expected cursor at origin, got %+v", got)
	}
}

func TestEditorDeleteSelectionSingleLine(t *testing.T) {
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

	e.DeleteSelection()

	if got := e.Line(0); got != "world" {
		t.Fatalf("expected %q, got %q", "world", got)
	}

	if got := e.Cursor(); got != (Position{Line: 0, Column: 0}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{}, got)
	}
}

func TestEditorCutSelection(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello world")

	e := NewEditor(b)
	e.SelectAll()

	text, ok := e.CutSelection()
	if !ok {
		t.Fatal("expected cut to return selected text")
	}

	if text != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", text)
	}

	if got := e.Line(0); got != "" {
		t.Fatalf("expected empty line, got %q", got)
	}
}

func TestEditorInsertText(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")

	e := NewEditor(b)
	e.cursor.Column = e.buffer.LineLength(0)

	e.InsertText(" world")

	if got := e.Line(0); got != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", got)
	}

	if got := e.Cursor(); got != (Position{Line: 0, Column: 11}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 0, Column: 11}, got)
	}
}

func TestEditorInsertMultilineText(t *testing.T) {
	b := buffer.NewBuffer()
	insertText(b, 0, "hello")

	e := NewEditor(b)
	e.cursor.Column = e.buffer.LineLength(0)

	e.InsertText(" world\nsecond line\nthird line")

	if got := e.LineCount(); got != 3 {
		t.Fatalf("expected 3 lines, got %d", got)
	}

	if got := e.Line(0); got != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", got)
	}

	if got := e.Line(1); got != "second line" {
		t.Fatalf("expected %q, got %q", "second line", got)
	}

	if got := e.Line(2); got != "third line" {
		t.Fatalf("expected %q, got %q", "third line", got)
	}

	if got := e.Cursor(); got != (Position{Line: 2, Column: 10}) {
		t.Fatalf("expected cursor %+v, got %+v", Position{Line: 2, Column: 10}, got)
	}
}
