package fileio

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	const original = "hello\nbeautiful\nworld"

	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	b, err := Open(path)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	defer func() {
		if err := b.Close(); err != nil {
			t.Fatalf("Close failed: %v", err)
		}
	}()

	if got := b.LineCount(); got != 3 {
		t.Fatalf("expected 3 lines, got %d", got)
	}

	if got := b.Line(1); got != "beautiful" {
		t.Fatalf("expected %q, got %q", "beautiful", got)
	}
}

func TestOpenEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	b, err := Open(path)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	defer func() {
		if err := b.Close(); err != nil {
			t.Fatalf("Close failed: %v", err)
		}
	}()

	if b.LineCount() != 1 {
		t.Fatalf("expected one empty line, got %d", b.LineCount())
	}
}

func TestOpenOrCreateCreatesMissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new.txt")

	b, err := OpenOrCreate(path)
	if err != nil {
		t.Fatalf("OpenOrCreate failed: %v", err)
	}

	defer func() {
		if err := b.Close(); err != nil {
			t.Fatalf("Close failed: %v", err)
		}
	}()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat created file: %v", err)
	}

	if info.Size() != 0 {
		t.Fatalf("expected empty file, got %d bytes", info.Size())
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	b, err := OpenOrCreate(path)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := b.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	b.InsertText(0, 0, "hello\nworld")

	if err := Save(path, b); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if got := string(data); got != "hello\nworld" {
		t.Fatalf("expected %q, got %q", "hello\nworld", got)
	}
}

func TestSaveReplacesExistingFileAtomically(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if err := os.WriteFile(
		path,
		[]byte("old"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	b, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := b.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	b.ReplaceRange(0, 0, 0, 3, "new")

	if err := Save(path, b); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if got := string(data); got != "new" {
		t.Fatalf("expected %q, got %q", "new", got)
	}
}
