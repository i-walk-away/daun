package editor

import "testing"

func TestEditorMoveLeft(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	editor := NewEditor(buffer)
	editor.cursor.Column = 3

	editor.MoveLeft()

	if editor.cursor.Column != 2 {
		t.Errorf("expected cursor column %d, got %d", 2, editor.cursor.Column)
	}
}

func TestEditorMoveLeftAtLineStart(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")
	buffer.InsertLine(1, "world")

	editor := NewEditor(buffer)
	editor.cursor.Line = 1
	editor.cursor.Column = 0

	editor.MoveLeft()

	if editor.cursor.Line != 0 {
		t.Errorf("expected cursor line %d, got %d", 0, editor.cursor.Line)
	}

	if editor.cursor.Column != 5 {
		t.Errorf("expected cursor column %d, got %d", 5, editor.cursor.Column)
	}
}

func TestEditorMoveLeftAtBufferStart(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	editor := NewEditor(buffer)

	editor.MoveLeft()

	if editor.cursor.Line != 0 || editor.cursor.Column != 0 {
		t.Errorf(
			"expected cursor at (0, 0), got (%d, %d)",
			editor.cursor.Line,
			editor.cursor.Column,
		)
	}
}

func TestEditorMoveRight(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	editor := NewEditor(buffer)

	editor.MoveRight()

	if editor.cursor.Column != 1 {
		t.Errorf("expected cursor column %d, got %d", 1, editor.cursor.Column)
	}
}

func TestEditorMoveRightAtLineEnd(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")
	buffer.InsertLine(1, "world")

	editor := NewEditor(buffer)
	editor.cursor.Column = 5

	editor.MoveRight()

	if editor.cursor.Line != 1 {
		t.Errorf("expected cursor line %d, got %d", 1, editor.cursor.Line)
	}

	if editor.cursor.Column != 0 {
		t.Errorf("expected cursor column %d, got %d", 0, editor.cursor.Column)
	}
}

func TestEditorMoveRightAtBufferEnd(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	editor := NewEditor(buffer)
	editor.cursor.Column = 5

	editor.MoveRight()

	if editor.cursor.Line != 0 || editor.cursor.Column != 5 {
		t.Errorf(
			"expected cursor at (0, 5), got (%d, %d)",
			editor.cursor.Line,
			editor.cursor.Column,
		)
	}
}

func TestEditorMoveUp(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")
	buffer.InsertLine(1, "world")

	editor := NewEditor(buffer)
	editor.cursor.Line = 1
	editor.cursor.Column = 3

	editor.MoveUp()

	if editor.cursor.Line != 0 {
		t.Errorf("expected cursor line %d, got %d", 0, editor.cursor.Line)
	}

	if editor.cursor.Column != 3 {
		t.Errorf("expected cursor column %d, got %d", 3, editor.cursor.Column)
	}
}

func TestEditorMoveUpToShorterLine(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hi")
	buffer.InsertLine(1, "hello")

	editor := NewEditor(buffer)
	editor.cursor.Line = 1
	editor.cursor.Column = 5

	editor.MoveUp()

	if editor.cursor.Line != 0 {
		t.Errorf("expected cursor line %d, got %d", 0, editor.cursor.Line)
	}

	if editor.cursor.Column != 2 {
		t.Errorf("expected cursor column %d, got %d", 2, editor.cursor.Column)
	}
}

func TestEditorMoveUpAtFirstLine(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	editor := NewEditor(buffer)
	editor.cursor.Column = 3

	editor.MoveUp()

	if editor.cursor.Line != 0 {
		t.Errorf("expected cursor line %d, got %d", 0, editor.cursor.Line)
	}

	if editor.cursor.Column != 0 {
		t.Errorf("expected cursor column %d, got %d", 0, editor.cursor.Column)
	}
}

func TestEditorMoveDown(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")
	buffer.InsertLine(1, "world")

	editor := NewEditor(buffer)
	editor.cursor.Column = 3

	editor.MoveDown()

	if editor.cursor.Line != 1 {
		t.Errorf("expected cursor line %d, got %d", 1, editor.cursor.Line)
	}

	if editor.cursor.Column != 3 {
		t.Errorf("expected cursor column %d, got %d", 3, editor.cursor.Column)
	}
}

func TestEditorMoveDownToShorterLine(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")
	buffer.InsertLine(1, "hi")

	editor := NewEditor(buffer)
	editor.cursor.Column = 5

	editor.MoveDown()

	if editor.cursor.Line != 1 {
		t.Errorf("expected cursor line %d, got %d", 1, editor.cursor.Line)
	}

	if editor.cursor.Column != 2 {
		t.Errorf("expected cursor column %d, got %d", 2, editor.cursor.Column)
	}
}

func TestEditorMoveDownAtLastLine(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	editor := NewEditor(buffer)
	editor.cursor.Column = 2

	editor.MoveDown()

	if editor.cursor.Line != 0 {
		t.Errorf("expected cursor line %d, got %d", 0, editor.cursor.Line)
	}

	if editor.cursor.Column != 5 {
		t.Errorf("expected cursor column %d, got %d", 5, editor.cursor.Column)
	}
}
