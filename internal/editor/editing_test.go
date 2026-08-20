package editor

import "testing"

func TestEditorInsert(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	editor := NewEditor(buffer)
	editor.cursor.Column = 2

	editor.Insert('X')

	if got := buffer.Line(0); got != "heXllo" {
		t.Errorf("expected line %q, got %q", "heXllo", got)
	}

	if editor.cursor.Column != 3 {
		t.Errorf("expected cursor column %d, got %d", 3, editor.cursor.Column)
	}
}

func TestEditorInsertUnicode(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "привет")

	editor := NewEditor(buffer)
	editor.cursor.Column = 3

	editor.Insert('🚀')

	if got := buffer.Line(0); got != "при🚀вет" {
		t.Errorf("expected line %q, got %q", "при🚀вет", got)
	}

	if editor.cursor.Column != 4 {
		t.Errorf("expected cursor column %d, got %d", 4, editor.cursor.Column)
	}
}

func TestEditorBackspace(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	editor := NewEditor(buffer)
	editor.cursor.Column = 3

	editor.Backspace()

	if got := buffer.Line(0); got != "helo" {
		t.Errorf("expected line %q, got %q", "helo", got)
	}

	if editor.cursor.Column != 2 {
		t.Errorf("expected cursor column %d, got %d", 2, editor.cursor.Column)
	}
}

func TestEditorBackspaceAtLineStart(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")
	buffer.InsertLine(1, "world")

	editor := NewEditor(buffer)
	editor.cursor.Line = 1
	editor.cursor.Column = 0

	editor.Backspace()

	if buffer.LineCount() != 1 {
		t.Errorf("expected %d line, got %d", 1, buffer.LineCount())
	}

	if got := buffer.Line(0); got != "helloworld" {
		t.Errorf("expected line %q, got %q", "helloworld", got)
	}

	if editor.cursor.Line != 0 {
		t.Errorf("expected cursor line %d, got %d", 0, editor.cursor.Line)
	}

	if editor.cursor.Column != 10 {
		t.Errorf("expected cursor column %d, got %d", 10, editor.cursor.Column)
	}
}

func TestEditorBackspaceAtBufferStart(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	editor := NewEditor(buffer)

	editor.Backspace()

	if got := buffer.Line(0); got != "hello" {
		t.Errorf("expected line %q, got %q", "hello", got)
	}

	if editor.cursor.Line != 0 || editor.cursor.Column != 0 {
		t.Errorf(
			"expected cursor at (0, 0), got (%d, %d)",
			editor.cursor.Line,
			editor.cursor.Column,
		)
	}
}

func TestEditorEnter(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello world")

	editor := NewEditor(buffer)
	editor.cursor.Column = 5

	editor.Enter()

	if buffer.LineCount() != 2 {
		t.Errorf("expected %d lines, got %d", 2, buffer.LineCount())
	}

	if got := buffer.Line(0); got != "hello" {
		t.Errorf("expected first line %q, got %q", "hello", got)
	}

	if got := buffer.Line(1); got != " world" {
		t.Errorf("expected second line %q, got %q", " world", got)
	}

	if editor.cursor.Line != 1 {
		t.Errorf("expected cursor line %d, got %d", 1, editor.cursor.Line)
	}

	if editor.cursor.Column != 0 {
		t.Errorf("expected cursor column %d, got %d", 0, editor.cursor.Column)
	}
}

func TestEditorEnterAtLineStart(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	editor := NewEditor(buffer)

	editor.Enter()

	if buffer.LineCount() != 2 {
		t.Errorf("expected %d lines, got %d", 2, buffer.LineCount())
	}

	if got := buffer.Line(0); got != "" {
		t.Errorf("expected first line %q, got %q", "", got)
	}

	if got := buffer.Line(1); got != "hello" {
		t.Errorf("expected second line %q, got %q", "hello", got)
	}

	if editor.cursor.Line != 1 || editor.cursor.Column != 0 {
		t.Errorf(
			"expected cursor at (1, 0), got (%d, %d)",
			editor.cursor.Line,
			editor.cursor.Column,
		)
	}
}

func TestEditorEnterAtLineEnd(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	editor := NewEditor(buffer)
	editor.cursor.Column = 5

	editor.Enter()

	if buffer.LineCount() != 2 {
		t.Errorf("expected %d lines, got %d", 2, buffer.LineCount())
	}

	if got := buffer.Line(0); got != "hello" {
		t.Errorf("expected first line %q, got %q", "hello", got)
	}

	if got := buffer.Line(1); got != "" {
		t.Errorf("expected second line %q, got %q", "", got)
	}

	if editor.cursor.Line != 1 || editor.cursor.Column != 0 {
		t.Errorf(
			"expected cursor at (1, 0), got (%d, %d)",
			1,
			0,
		)
	}
}
