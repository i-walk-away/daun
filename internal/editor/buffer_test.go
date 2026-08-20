package editor

import "testing"

func TestNewBuffer(t *testing.T) {
	buffer := NewBuffer()

	if buffer.LineCount() != 1 {
		t.Errorf("expected 1 line, got %d", buffer.LineCount())
	}

	if got := buffer.Line(0); got != "" {
		t.Errorf("expected empty line, got %q", got)
	}
}

func TestBufferLine(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	if got := buffer.Line(0); got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestBufferSetLine(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	buffer.SetLine(0, "world")

	if got := buffer.Line(0); got != "world" {
		t.Errorf("expected %q, got %q", "world", got)
	}
}

func TestBufferInsertLine(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	buffer.InsertLine(1, "world")

	if buffer.LineCount() != 2 {
		t.Errorf("expected 2 lines, got %d", buffer.LineCount())
	}

	if got := buffer.Line(0); got != "hello" {
		t.Errorf("expected first line %q, got %q", "hello", got)
	}

	if got := buffer.Line(1); got != "world" {
		t.Errorf("expected second line %q, got %q", "world", got)
	}
}

func TestBufferInsertLineInMiddle(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")
	buffer.InsertLine(1, "world")
	buffer.InsertLine(2, "!!!")

	buffer.InsertLine(1, "new")

	expected := []string{
		"hello",
		"new",
		"world",
		"!!!",
	}

	if buffer.LineCount() != len(expected) {
		t.Errorf("expected %d lines, got %d", len(expected), buffer.LineCount())
	}

	for i, want := range expected {
		if got := buffer.Line(i); got != want {
			t.Errorf("line %d: expected %q, got %q", i, want, got)
		}
	}
}

func TestBufferInsertLineAtBeginning(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "world")

	buffer.InsertLine(0, "hello")

	if buffer.LineCount() != 2 {
		t.Errorf("expected 2 lines, got %d", buffer.LineCount())
	}

	if got := buffer.Line(0); got != "hello" {
		t.Errorf("expected first line %q, got %q", "hello", got)
	}

	if got := buffer.Line(1); got != "world" {
		t.Errorf("expected second line %q, got %q", "world", got)
	}
}

func TestBufferDeleteLine(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")
	buffer.InsertLine(1, "world")

	buffer.DeleteLine(1)

	if buffer.LineCount() != 1 {
		t.Errorf("expected 1 line, got %d", buffer.LineCount())
	}

	if got := buffer.Line(0); got != "hello" {
		t.Errorf("expected line %q, got %q", "hello", got)
	}
}

func TestBufferDeleteLineInMiddle(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")
	buffer.InsertLine(1, "world")
	buffer.InsertLine(2, "!!!")

	buffer.DeleteLine(1)

	expected := []string{
		"hello",
		"!!!",
	}

	if buffer.LineCount() != len(expected) {
		t.Errorf("expected %d lines, got %d", len(expected), buffer.LineCount())
	}

	for i, want := range expected {
		if got := buffer.Line(i); got != want {
			t.Errorf("line %d: expected %q, got %q", i, want, got)
		}
	}
}

func TestBufferLineCount(t *testing.T) {
	buffer := NewBuffer()

	if got := buffer.LineCount(); got != 1 {
		t.Errorf("expected 1 line, got %d", got)
	}

	buffer.InsertLine(1, "hello")

	if got := buffer.LineCount(); got != 2 {
		t.Errorf("expected 2 lines, got %d", got)
	}

	buffer.DeleteLine(1)

	if got := buffer.LineCount(); got != 1 {
		t.Errorf("expected 1 line, got %d", got)
	}
}

func TestBufferLineLength(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello")

	if got := buffer.LineLength(0); got != 5 {
		t.Errorf("expected length %d, got %d", 5, got)
	}
}

func TestBufferLineLengthUnicode(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "привет")

	if got := buffer.LineLength(0); got != 6 {
		t.Errorf("expected length %d, got %d", 6, got)
	}
}

func TestBufferLineLengthEmoji(t *testing.T) {
	buffer := NewBuffer()
	buffer.SetLine(0, "hello 🚀")

	if got := buffer.LineLength(0); got != 7 {
		t.Errorf("expected length %d, got %d", 7, got)
	}
}
