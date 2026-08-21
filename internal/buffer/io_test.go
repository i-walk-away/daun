package buffer

import (
	"bytes"
	"strings"
	"testing"
)

func TestBufferWriteTo(t *testing.T) {
	b := NewBufferFromText("hello\nbeautiful\nworld")

	var dst bytes.Buffer

	n, err := b.WriteTo(&dst)
	if err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}

	if got := dst.String(); got != "hello\nbeautiful\nworld" {
		t.Fatalf("expected %q, got %q", "hello\nbeautiful\nworld", got)
	}

	if n != int64(dst.Len()) {
		t.Fatalf("expected %d bytes written, got %d", dst.Len(), n)
	}
}

func TestBufferWriteToAfterEdits(t *testing.T) {
	b := NewBufferFromText("hello world")

	b.Insert(0, 5, ',')
	b.Delete(0, 7)
	b.InsertText(0, b.LineLength(0), "!")

	var dst bytes.Buffer

	if _, err := b.WriteTo(&dst); err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}

	if got := dst.String(); got != "hello,world!" {
		t.Fatalf("expected %q, got %q", "hello,world!", got)
	}
}

func TestBufferWriteToLargeDocument(t *testing.T) {
	text := strings.Repeat("hello world\n", 100_000)
	b := NewBufferFromText(text)

	var dst bytes.Buffer

	n, err := b.WriteTo(&dst)
	if err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}

	if n != int64(len(text)) {
		t.Fatalf("expected %d bytes, got %d", len(text), n)
	}

	if dst.String() != text {
		t.Fatal("written document differs from source")
	}
}
