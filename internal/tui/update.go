package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/i-walk-away/daun/internal/editor"
)

// Update handles incoming terminal events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateViewport()

	case tea.MouseWheelMsg:
		switch msg.Button {
		case tea.MouseWheelUp:
			m.viewport.ScrollUp(3)

		case tea.MouseWheelDown:
			m.viewport.ScrollDown(3, m.editor.LineCount())
		}

	case tea.KeyPressMsg:
		cmd, quit := m.handleKey(msg)

		m.updateViewport()

		if quit {
			return m, tea.Quit
		}

		return m, cmd
	}

	return m, nil
}

// handleKey handles keyboard input.
//
// It returns a command when the key requires an asynchronous Bubble Tea
// operation and returns true when the application should quit.
func (m *Model) handleKey(msg tea.KeyPressMsg) (tea.Cmd, bool) {
	shift := msg.Mod.Contains(tea.ModShift)
	ctrl := msg.Mod.Contains(tea.ModCtrl)

	switch {
	case ctrl && msg.Code == 'a':
		m.editor.SelectAll()

	case ctrl && msg.Code == 'c':
		if text, ok := m.editor.SelectedText(); ok {
			if err := m.clipboard.Copy(text); err != nil {
				return nil, false
			}
		}

	case ctrl && msg.Code == 'x':
		if text, ok := m.editor.CutSelection(); ok {
			if err := m.clipboard.Copy(text); err != nil {
				return nil, false
			}
		}

	case ctrl && msg.Code == 'v':
		text, err := m.clipboard.Paste()
		if err != nil {
			return nil, false
		}

		m.editor.InsertText(text)

	case ctrl && msg.Code == 'z' && shift:
		m.editor.Redo()

	case ctrl && msg.Code == 'z':
		m.editor.Undo()

	case ctrl && msg.Code == 'y':
		m.editor.Redo()

	case msg.Code == tea.KeyEscape:
		return nil, true

	case ctrl && msg.Code == tea.KeyLeft:
		m.editor.Move(
			editor.DirectionLeft,
			editor.MoveByWord,
			shift,
		)

	case ctrl && msg.Code == tea.KeyRight:
		m.editor.Move(
			editor.DirectionRight,
			editor.MoveByWord,
			shift,
		)

	case ctrl && msg.Code == 'h':
		m.editor.DeleteWordBackwards()

	case msg.Code == tea.KeyBackspace:
		m.editor.Backspace()

	case msg.Code == tea.KeyEnter:
		m.editor.Enter()

	case msg.Code == tea.KeyLeft:
		m.editor.Move(
			editor.DirectionLeft,
			editor.MoveByCharacter,
			shift,
		)

	case msg.Code == tea.KeyRight:
		m.editor.Move(
			editor.DirectionRight,
			editor.MoveByCharacter,
			shift,
		)

	case msg.Code == tea.KeyUp:
		m.editor.Move(
			editor.DirectionUp,
			editor.MoveByCharacter,
			shift,
		)

	case msg.Code == tea.KeyDown:
		m.editor.Move(
			editor.DirectionDown,
			editor.MoveByCharacter,
			shift,
		)

	case msg.Code == tea.KeySpace:
		m.editor.Insert(' ')

	default:
		if msg.Text != "" {
			m.editor.InsertText(msg.Text)
		}
	}

	return nil, false
}

func (m *Model) updateViewport() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	lineNumberWidth := lineNumberWidth(m.editor.LineCount())
	contentHeight := max(m.height-1, 0)

	m.viewport.SetSize(m.width, contentHeight)

	cursor := m.editor.Cursor()

	m.viewport.EnsureCursorVisible(
		cursor,
		m.editor.Line(cursor.Line),
		m.editor.LineCount(),
		lineNumberWidth,
	)
}
