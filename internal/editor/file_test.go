package editor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/i-walk-away/daun/internal/fileio"
)

func TestEditorSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	b, err := fileio.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	e := NewFileEditor(b, path)
	e.Move(
		DirectionRight,
		MoveByCharacter,
		false,
	)
	e.Move(
		DirectionRight,
		MoveByCharacter,
		false,
	)
	e.Move(
		DirectionRight,
		MoveByCharacter,
		false,
	)
	e.Move(
		DirectionRight,
		MoveByCharacter,
		false,
	)
	e.Move(
		DirectionRight,
		MoveByCharacter,
		false,
	)

	e.InsertText(" world")

	if !e.Modified() {
		t.Fatal("expected modified editor")
	}

	if err := e.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if e.Modified() {
		t.Fatal("expected editor to be clean after save")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if got := string(data); got != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", got)
	}

	if err := e.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestEditorUndoRestoresSavedRevision(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")

	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	b, err := fileio.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	e := NewFileEditor(b, path)
	defer func() {
		if err := e.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	if e.Modified() {
		t.Fatal("newly opened file must not be modified")
	}

	for range 5 {
		e.Move(
			DirectionRight,
			MoveByCharacter,
			false,
		)
	}

	e.Insert('!')

	if !e.Modified() {
		t.Fatal("expected modified editor after insert")
	}

	e.Undo()

	if e.Modified() {
		t.Fatal("expected undo to restore saved revision")
	}

	if got := e.Line(0); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}
