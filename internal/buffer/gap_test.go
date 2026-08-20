package buffer

import "testing"

func TestGapBufferInsert(t *testing.T) {
	var g gapBuffer

	g.Insert('h')
	g.Insert('e')
	g.Insert('l')
	g.Insert('l')
	g.Insert('o')

	if got := g.String(); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}

	if got := g.Len(); got != 5 {
		t.Fatalf("expected length 5, got %d", got)
	}
}

func TestGapBufferInsertAtPosition(t *testing.T) {
	var g gapBuffer

	for _, r := range "helo" {
		g.Insert(r)
	}

	g.MoveGap(3)
	g.Insert('l')

	if got := g.String(); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestGapBufferMoveGap(t *testing.T) {
	var g gapBuffer

	for _, r := range "hello" {
		g.Insert(r)
	}

	g.MoveGap(0)

	if got := g.String(); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}

	g.MoveGap(5)

	if got := g.String(); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}

	g.MoveGap(2)

	if got := g.String(); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestGapBufferDelete(t *testing.T) {
	var g gapBuffer

	for _, r := range "hello" {
		g.Insert(r)
	}

	g.MoveGap(5)
	g.Delete()

	if got := g.String(); got != "hell" {
		t.Fatalf("expected %q, got %q", "hell", got)
	}
}

func TestGapBufferDeleteRange(t *testing.T) {
	var g gapBuffer

	for _, r := range "hello beautiful world" {
		g.Insert(r)
	}

	g.DeleteRange(6, 15)

	if got := g.String(); got != "hello  world" {
		t.Fatalf("expected %q, got %q", "hello  world", got)
	}
}

func TestGapBufferDeleteRangeAtBeginning(t *testing.T) {
	var g gapBuffer

	for _, r := range "hello world" {
		g.Insert(r)
	}

	g.DeleteRange(0, 6)

	if got := g.String(); got != "world" {
		t.Fatalf("expected %q, got %q", "world", got)
	}
}

func TestGapBufferDeleteRangeAtEnd(t *testing.T) {
	var g gapBuffer

	for _, r := range "hello world" {
		g.Insert(r)
	}

	g.DeleteRange(5, 11)

	if got := g.String(); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestGapBufferDeleteRangeEmpty(t *testing.T) {
	var g gapBuffer

	for _, r := range "hello" {
		g.Insert(r)
	}

	g.DeleteRange(2, 2)

	if got := g.String(); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestGapBufferUnicode(t *testing.T) {
	var g gapBuffer

	for _, r := range "привет 👋" {
		g.Insert(r)
	}

	if got := g.String(); got != "привет 👋" {
		t.Fatalf("expected Unicode string, got %q", got)
	}

	if got := g.Len(); got != 8 {
		t.Fatalf("expected rune length 8, got %d", got)
	}
}
