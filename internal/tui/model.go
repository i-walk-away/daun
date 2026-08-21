package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/i-walk-away/daun/internal/editor"
)

// Model represents the state of the editor TUI.
type Model struct {
	editor   *editor.Editor
	viewport Viewport
	width    int
	height   int
}

// NewModel creates a new TUI model from an editor.
func NewModel(editor *editor.Editor) Model {
	return Model{
		editor: editor,
	}
}

// Init initializes the TUI.
func (m Model) Init() tea.Cmd {
	return nil
}
