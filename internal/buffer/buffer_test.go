package buffer

import "testing"

func insertText(b *Buffer, line int, text string) {
	for _, r := range text {
		b.Insert(line, b.LineLength(line), r)
	}
}

func TestNewBuffer(t *testing.T) {
	b := NewBuffer()

	if b.LineCount() != 1 {
		t.Fatalf("expected 1 line, got %d", b.LineCount())
	}

	if got := b.Line(0); got != "" {
		t.Fatalf("expected empty line, got %q", got)
	}
}

func TestBufferInsert(t *testing.T) {
	b := NewBuffer()

	b.Insert(0, 0, 'h')
	b.Insert(0, 1, 'e')
	b.Insert(0, 2, 'l')
	b.Insert(0, 3, 'l')
	b.Insert(0, 4, 'o')

	if got := b.Line(0); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}

	if got := b.LineLength(0); got != 5 {
		t.Fatalf("expected line length 5, got %d", got)
	}
}

func TestBufferInsertInMiddle(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "helo")

	b.Insert(0, 2, 'l')

	if got := b.Line(0); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestBufferInsertUnicode(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "hello ")

	b.Insert(0, b.LineLength(0), '世')
	b.Insert(0, b.LineLength(0), '界')

	if got := b.Line(0); got != "hello 世界" {
		t.Fatalf("expected %q, got %q", "hello 世界", got)
	}

	if got := b.LineLength(0); got != 8 {
		t.Fatalf("expected line length 8, got %d", got)
	}
}

func TestBufferDelete(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "hello")

	// Delete the character immediately before column 5.
	b.Delete(0, 5)

	if got := b.Line(0); got != "hell" {
		t.Fatalf("expected %q, got %q", "hell", got)
	}
}

func TestBufferDeleteInMiddle(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "hello")

	// Delete the character immediately before column 3.
	b.Delete(0, 3)

	if got := b.Line(0); got != "helo" {
		t.Fatalf("expected %q, got %q", "helo", got)
	}
}

func TestBufferDeleteRange(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "hello beautiful world")

	b.DeleteRange(0, 6, 15)

	if got := b.Line(0); got != "hello  world" {
		t.Fatalf("expected %q, got %q", "hello  world", got)
	}
}

func TestBufferInsertLine(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "first")

	b.InsertLine(1)
	insertText(b, 1, "second")

	if got := b.LineCount(); got != 2 {
		t.Fatalf("expected 2 lines, got %d", got)
	}

	if got := b.Line(0); got != "first" {
		t.Fatalf("expected line 0 to be %q, got %q", "first", got)
	}

	if got := b.Line(1); got != "second" {
		t.Fatalf("expected line 1 to be %q, got %q", "second", got)
	}
}

func TestBufferInsertLineInMiddle(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "first")

	b.InsertLine(1)
	insertText(b, 1, "third")

	b.InsertLine(1)
	insertText(b, 1, "second")

	if got := b.LineCount(); got != 3 {
		t.Fatalf("expected 3 lines, got %d", got)
	}

	expected := []string{
		"first",
		"second",
		"third",
	}

	for i, want := range expected {
		if got := b.Line(i); got != want {
			t.Fatalf("expected line %d to be %q, got %q", i, want, got)
		}
	}
}

func TestBufferDeleteLine(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "first")

	b.InsertLine(1)
	insertText(b, 1, "second")

	b.DeleteLine(0)

	if got := b.LineCount(); got != 1 {
		t.Fatalf("expected 1 line, got %d", got)
	}

	if got := b.Line(0); got != "second" {
		t.Fatalf("expected %q, got %q", "second", got)
	}
}

func TestBufferDeleteMiddleLine(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "first")

	b.InsertLine(1)
	insertText(b, 1, "second")

	b.InsertLine(2)
	insertText(b, 2, "third")

	b.DeleteLine(1)

	if got := b.LineCount(); got != 2 {
		t.Fatalf("expected 2 lines, got %d", got)
	}

	expected := []string{
		"first",
		"third",
	}

	for i, want := range expected {
		if got := b.Line(i); got != want {
			t.Fatalf("expected line %d to be %q, got %q", i, want, got)
		}
	}
}

func TestBufferDeleteLastLine(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "first")

	b.InsertLine(1)
	insertText(b, 1, "second")

	b.DeleteLine(1)

	if got := b.LineCount(); got != 1 {
		t.Fatalf("expected 1 line, got %d", got)
	}

	if got := b.Line(0); got != "first" {
		t.Fatalf("expected %q, got %q", "first", got)
	}
}

func TestBufferSplitLine(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "hello world")

	b.SplitLine(0, 5)

	if got := b.LineCount(); got != 2 {
		t.Fatalf("expected 2 lines, got %d", got)
	}

	if got := b.Line(0); got != "hello" {
		t.Fatalf("expected line 0 to be %q, got %q", "hello", got)
	}

	if got := b.Line(1); got != " world" {
		t.Fatalf("expected line 1 to be %q, got %q", " world", got)
	}
}

func TestBufferSplitLineAtStart(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "hello")

	b.SplitLine(0, 0)

	if got := b.LineCount(); got != 2 {
		t.Fatalf("expected 2 lines, got %d", got)
	}

	if got := b.Line(0); got != "" {
		t.Fatalf("expected line 0 to be empty, got %q", got)
	}

	if got := b.Line(1); got != "hello" {
		t.Fatalf("expected line 1 to be %q, got %q", "hello", got)
	}
}

func TestBufferSplitLineAtEnd(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "hello")

	b.SplitLine(0, 5)

	if got := b.LineCount(); got != 2 {
		t.Fatalf("expected 2 lines, got %d", got)
	}

	if got := b.Line(0); got != "hello" {
		t.Fatalf("expected line 0 to be %q, got %q", "hello", got)
	}

	if got := b.Line(1); got != "" {
		t.Fatalf("expected line 1 to be empty, got %q", got)
	}
}

func TestBufferJoinLines(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "hello")

	b.InsertLine(1)
	insertText(b, 1, "world")

	b.JoinLines(0)

	if got := b.LineCount(); got != 1 {
		t.Fatalf("expected 1 line, got %d", got)
	}

	if got := b.Line(0); got != "helloworld" {
		t.Fatalf("expected %q, got %q", "helloworld", got)
	}
}

func TestBufferJoinEmptyLines(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "hello")
	b.InsertLine(1)

	b.JoinLines(0)

	if got := b.LineCount(); got != 1 {
		t.Fatalf("expected 1 line, got %d", got)
	}

	if got := b.Line(0); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestBufferLineLength(t *testing.T) {
	b := NewBuffer()

	insertText(b, 0, "hello")

	if got := b.LineLength(0); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}

	insertText(b, 0, " 世界")

	if got := b.LineLength(0); got != 8 {
		t.Fatalf("expected 8, got %d", got)
	}
}
