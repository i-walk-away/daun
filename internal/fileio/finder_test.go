package fileio

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileFinderSearch(t *testing.T) {
	root := t.TempDir()

	files := []string{
		"README.md",
		"src/main.go",
		"src/editor.go",
		"docs/README.md",
	}

	for _, file := range files {
		path := filepath.Join(root, file)

		if err := os.MkdirAll(
			filepath.Dir(path),
			0o755,
		); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(
			path,
			[]byte("test"),
			0o644,
		); err != nil {
			t.Fatal(err)
		}
	}

	finder := NewFileFinder(root)

	results, err := finder.Search(
		context.Background(),
		"readme",
		100,
	)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf(
			"expected 2 results, got %d",
			len(results),
		)
	}

	if results[0].DisplayPath != "README.md" {
		t.Fatalf(
			"expected README.md, got %q",
			results[0].DisplayPath,
		)
	}

	if results[1].DisplayPath != "docs/README.md" {
		t.Fatalf(
			"expected docs/README.md, got %q",
			results[1].DisplayPath,
		)
	}
}

func TestFileFinderSearchRespectsLimit(t *testing.T) {
	root := t.TempDir()

	for i := range 20 {
		path := filepath.Join(
			root,
			"file",
			string(rune('a'+i)),
			"target.txt",
		)

		if err := os.MkdirAll(
			filepath.Dir(path),
			0o755,
		); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(
			path,
			nil,
			0o644,
		); err != nil {
			t.Fatal(err)
		}
	}

	finder := NewFileFinder(root)

	results, err := finder.Search(
		context.Background(),
		"target",
		5,
	)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 5 {
		t.Fatalf(
			"expected 5 results, got %d",
			len(results),
		)
	}
}

func TestFileFinderSearchCanBeCancelled(t *testing.T) {
	root := t.TempDir()

	for i := range 100 {
		path := filepath.Join(
			root,
			"file",
			"target.txt",
			string(rune('a'+i%26)),
		)

		if err := os.MkdirAll(
			filepath.Dir(path),
			0o755,
		); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(
			path,
			nil,
			0o644,
		); err != nil {
			t.Fatal(err)
		}
	}

	ctx, cancel := context.WithCancel(
		context.Background(),
	)
	cancel()

	finder := NewFileFinder(root)

	_, err := finder.Search(
		ctx,
		"target",
		100,
	)

	if err != context.Canceled {
		t.Fatalf(
			"expected context.Canceled, got %v",
			err,
		)
	}
}
