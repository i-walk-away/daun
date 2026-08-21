package editor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/i-walk-away/daun/internal/fileio"
)

func TestEditorOpenFile(t *testing.T) {
	dir := t.TempDir()

	first := filepath.Join(dir, "first.txt")
	second := filepath.Join(dir, "second.txt")

	if err := os.WriteFile(
		first,
		[]byte("first"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(
		second,
		[]byte("second"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	b, err := fileio.Open(first)
	if err != nil {
		t.Fatal(err)
	}

	e := NewFileEditor(b, first)
	defer func() {
		if err := e.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if err := e.OpenFile(second); err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}

	if got := e.FilePath(); got != second {
		t.Fatalf(
			"expected path %q, got %q",
			second,
			got,
		)
	}

	if got := e.Line(0); got != "second" {
		t.Fatalf(
			"expected %q, got %q",
			"second",
			got,
		)
	}

	if e.Modified() {
		t.Fatal("opened file should not be modified")
	}
}

func TestEditorOpenFileRejectsUnsavedChanges(t *testing.T) {
	dir := t.TempDir()

	first := filepath.Join(dir, "first.txt")
	second := filepath.Join(dir, "second.txt")

	if err := os.WriteFile(
		first,
		[]byte("first"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(
		second,
		[]byte("second"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	b, err := fileio.Open(first)
	if err != nil {
		t.Fatal(err)
	}

	e := NewFileEditor(b, first)
	defer func() {
		if err := e.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	e.InsertText("!")

	if err := e.OpenFile(second); err == nil {
		t.Fatal("expected OpenFile to reject unsaved changes")
	}

	if got := e.FilePath(); got != first {
		t.Fatalf(
			"expected current file %q, got %q",
			first,
			got,
		)
	}
}
