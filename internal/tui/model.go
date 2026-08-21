package tui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/i-walk-away/daun/internal/clipboard"
	"github.com/i-walk-away/daun/internal/editor"
	"github.com/i-walk-away/daun/internal/fileio"
)

// Model represents the state of the editor TUI.
type Model struct {
	editor     *editor.Editor
	clipboard  clipboard.Clipboard
	fileFinder *fileio.FileFinder

	viewport Viewport

	width  int
	height int

	message *Message
	search  *SearchState
}

// NewModel creates a new TUI model from an editor, clipboard backend,
// and file finder.
func NewModel(
	editor *editor.Editor,
	clipboard clipboard.Clipboard,
	fileFinder *fileio.FileFinder,
) Model {
	input := textinput.New()

	input.Prompt = "/ "
	input.Placeholder = "Search text..."
	input.CharLimit = 4096
	input.SetVirtualCursor(false)

	return Model{
		editor:     editor,
		clipboard:  clipboard,
		fileFinder: fileFinder,
		search: &SearchState{
			Input: input,
		},
	}
}

// Init initializes the TUI.
func (m Model) Init() tea.Cmd {
	return nil
}
