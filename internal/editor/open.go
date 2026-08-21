package editor

import (
	"fmt"

	"github.com/i-walk-away/daun/internal/fileio"
)

// OpenFile replaces the current document with the file at path.
//
// The operation is rejected when the current document contains unsaved
// changes. The new document becomes the current file and starts with a
// clean revision.
func (e *Editor) OpenFile(path string) error {
	if e.Modified() {
		return fmt.Errorf("cannot open file: document has unsaved changes")
	}

	nextBuffer, err := fileio.Open(path)
	if err != nil {
		return fmt.Errorf("open %q: %w", path, err)
	}

	oldBuffer := e.buffer

	if err := oldBuffer.Close(); err != nil {
		_ = nextBuffer.Close()

		return fmt.Errorf("close current document: %w", err)
	}

	e.buffer = nextBuffer
	e.filePath = path

	e.cursor = Position{}
	e.selection = nil

	e.revision = 0
	e.savedRevision = 0

	e.undoStack = nil
	e.redoStack = nil

	return nil
}
