package buffer

import "testing"

func TestLineInsert(t *testing.T) {
	var line Line

	for _, r := range "hello" {
		line.Insert(r)
	}

	if got := line.String(); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}

	if got := line.Len(); got != 5 {
		t.Fatalf("expected length 5, got %d", got)
	}
}

func TestLineInsertAtPosition(t *testing.T) {
	var line Line

	for _, r := range "helo" {
		line.Insert(r)
	}

	line.SetPosition(3)
	line.Insert('l')

	if got := line.String(); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestLineDelete(t *testing.T) {
	var line Line

	for _, r := range "hello" {
		line.Insert(r)
	}

	line.SetPosition(5)
	line.Delete()

	if got := line.String(); got != "hell" {
		t.Fatalf("expected %q, got %q", "hell", got)
	}
}

func TestLineDeleteRange(t *testing.T) {
	var line Line

	for _, r := range "hello world" {
		line.Insert(r)
	}

	line.DeleteRange(5, 6)

	if got := line.String(); got != "helloworld" {
		t.Fatalf("expected %q, got %q", "helloworld", got)
	}
}

func TestLineLenUnicode(t *testing.T) {
	var line Line

	for _, r := range "привет" {
		line.Insert(r)
	}

	if got := line.Len(); got != 6 {
		t.Fatalf("expected length 6, got %d", got)
	}
}

func TestLineStringEmpty(t *testing.T) {
	var line Line

	if got := line.String(); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}
