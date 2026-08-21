package editor

import (
	"fmt"

	"github.com/i-walk-away/daun/internal/fileio"
)

// Save writes the current document to its associated file.
//
// Saving does not change the editor's undo history. If saving succeeds,
// the current revision becomes the saved revision.
func (e *Editor) Save() error {
	if e.filePath == "" {
		return fmt.Errorf("editor has no associated file")
	}

	if err := fileio.Save(e.filePath, e.buffer); err != nil {
		return err
	}

	e.savedRevision = e.revision

	return nil
}

// SaveAs writes the current document to path and associates the editor
// with that file.
//
// The document remains at the same revision; only the associated path and
// saved revision are changed after a successful save.
func (e *Editor) SaveAs(path string) error {
	if path == "" {
		return fmt.Errorf("file path is empty")
	}

	if err := fileio.Save(path, e.buffer); err != nil {
		return err
	}

	e.filePath = path
	e.savedRevision = e.revision

	return nil
}
