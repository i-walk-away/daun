package editor

import (
	"strings"
)

type historyEntry struct {
	start        Position
	deletedText  string
	insertedText string

	beforeCursor    Position
	afterCursor     Position
	beforeSelection *Selection
	afterSelection  *Selection

	beforeRevision uint64
	afterRevision  uint64
}

// Undo undoes the most recent editing operation.
//
// If there are no operations to undo, Undo does nothing.
func (e *Editor) Undo() {
	if len(e.undoStack) == 0 {
		return
	}

	last := len(e.undoStack) - 1
	entry := e.undoStack[last]

	e.undoStack = e.undoStack[:last]

	end := advancePosition(entry.start, entry.insertedText)

	e.buffer.ReplaceRange(
		entry.start.Line,
		entry.start.Column,
		end.Line,
		end.Column,
		entry.deletedText,
	)

	e.cursor = entry.beforeCursor
	e.selection = cloneSelection(entry.beforeSelection)
	e.revision = entry.beforeRevision

	e.redoStack = append(e.redoStack, entry)
}

// Redo reapplies the most recently undone editing operation.
//
// If there are no operations to redo, Redo does nothing.
func (e *Editor) Redo() {
	if len(e.redoStack) == 0 {
		return
	}

	last := len(e.redoStack) - 1
	entry := e.redoStack[last]

	e.redoStack = e.redoStack[:last]

	end := advancePosition(entry.start, entry.deletedText)

	e.buffer.ReplaceRange(
		entry.start.Line,
		entry.start.Column,
		end.Line,
		end.Column,
		entry.insertedText,
	)

	e.cursor = entry.afterCursor
	e.selection = cloneSelection(entry.afterSelection)
	e.revision = entry.afterRevision

	e.undoStack = append(e.undoStack, entry)
}

func (e *Editor) applyEdit(
	start Position,
	end Position,
	replacement string,
	afterCursor Position,
	afterSelection *Selection,
) {
	beforeCursor := e.cursor
	beforeSelection := cloneSelection(e.selection)
	beforeRevision := e.revision

	deletedText := e.buffer.TextRange(
		start.Line,
		start.Column,
		end.Line,
		end.Column,
	)

	if deletedText == replacement &&
		beforeCursor == afterCursor &&
		sameSelection(beforeSelection, afterSelection) {
		return
	}

	e.buffer.ReplaceRange(
		start.Line,
		start.Column,
		end.Line,
		end.Column,
		replacement,
	)

	e.revision++

	e.cursor = afterCursor
	e.selection = cloneSelection(afterSelection)

	e.undoStack = append(e.undoStack, historyEntry{
		start:           start,
		deletedText:     deletedText,
		insertedText:    replacement,
		beforeCursor:    beforeCursor,
		afterCursor:     e.cursor,
		beforeSelection: beforeSelection,
		afterSelection:  cloneSelection(e.selection),
		beforeRevision:  beforeRevision,
		afterRevision:   e.revision,
	})

	e.redoStack = nil
}

func advancePosition(start Position, text string) Position {
	lines := strings.Split(text, "\n")

	if len(lines) == 1 {
		return Position{
			Line:   start.Line,
			Column: start.Column + len([]rune(text)),
		}
	}

	return Position{
		Line:   start.Line + len(lines) - 1,
		Column: len([]rune(lines[len(lines)-1])),
	}
}

func cloneSelection(selection *Selection) *Selection {
	if selection == nil {
		return nil
	}

	return &Selection{
		Anchor: selection.Anchor,
	}
}

func sameSelection(a, b *Selection) bool {
	if a == nil || b == nil {
		return a == b
	}

	return a.Anchor == b.Anchor
}
