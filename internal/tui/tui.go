package tui

import (
	"strings"

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

// Update handles incoming terminal events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyEscape:
			return m, tea.Quit

		case tea.KeyBackspace:
			m.editor.Backspace()

		case tea.KeyEnter:
			m.editor.Enter()

		case tea.KeyLeft:
			m.editor.MoveLeft()

		case tea.KeyRight:
			m.editor.MoveRight()

		case tea.KeyUp:
			m.editor.MoveUp()

		case tea.KeyDown:
			m.editor.MoveDown()

		case tea.KeySpace:
			m.editor.Insert(' ')

		default:
			if msg.Text != "" {
				for _, r := range msg.Text {
					m.editor.Insert(r)
				}
			}
		}
	}

	return m, nil
}

// View renders the current editor state.
func (m Model) View() tea.View {
	var b strings.Builder

	for i := 0; i < m.editor.LineCount(); i++ {
		b.WriteString(m.editor.Line(i))
		b.WriteRune('\n')
	}

	view := tea.NewView(b.String())

	cursor := m.editor.Cursor()

	view.Cursor = &tea.Cursor{
		X:     cursor.Column,
		Y:     cursor.Line,
		Shape: 2,
		Blink: true,
	}

	return view
}
