package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/i-walk-away/daun/internal/clipboard"
	"github.com/i-walk-away/daun/internal/editor"
)

// Model represents the state of the editor TUI.
type Model struct {
	editor    *editor.Editor
	clipboard clipboard.Clipboard

	viewport Viewport

	width  int
	height int

	message *Message
}

// NewModel creates a new TUI model from an editor and clipboard backend.
func NewModel(
	editor *editor.Editor,
	clipboard clipboard.Clipboard,
) Model {
	return Model{
		editor:    editor,
		clipboard: clipboard,
	}
}

// Init initializes the TUI.
func (m Model) Init() tea.Cmd {
	return nil
}
