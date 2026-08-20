package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/i-walk-away/daun/internal/editor"
)

// Model represents the TUI state.
type Model struct {
	editor *editor.Editor
}

// NewModel creates a new TUI model from an editor.
func NewModel(editor *editor.Editor) Model {
	return Model{
		editor: editor,
	}
}

// Init initializes the TUI model.
func (m Model) Init() tea.Cmd {
	return nil
}
