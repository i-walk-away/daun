package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/i-walk-away/daun/internal/editor"
)

// Update handles incoming terminal events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case msg.Code == tea.KeyEscape:
			return m, tea.Quit

		case msg.Mod.Contains(tea.ModCtrl) && msg.Code == tea.KeyLeft:
			m.editor.Move(
				editor.DirectionLeft,
				editor.MoveByWord,
				msg.Mod.Contains(tea.ModShift),
			)

		case msg.Mod.Contains(tea.ModCtrl) && msg.Code == tea.KeyRight:
			m.editor.Move(
				editor.DirectionRight,
				editor.MoveByWord,
				msg.Mod.Contains(tea.ModShift),
			)

		case msg.Mod.Contains(tea.ModCtrl) && msg.Code == 'h':
			m.editor.DeleteWordBackwards()

		case msg.Code == tea.KeyBackspace:
			m.editor.Backspace()

		case msg.Code == tea.KeyEnter:
			m.editor.Enter()

		case msg.Code == tea.KeyLeft:
			m.editor.Move(
				editor.DirectionLeft,
				editor.MoveByCharacter,
				msg.Mod.Contains(tea.ModShift),
			)

		case msg.Code == tea.KeyRight:
			m.editor.Move(
				editor.DirectionRight,
				editor.MoveByCharacter,
				msg.Mod.Contains(tea.ModShift),
			)

		case msg.Code == tea.KeyUp:
			m.editor.Move(
				editor.DirectionUp,
				editor.MoveByCharacter,
				msg.Mod.Contains(tea.ModShift),
			)

		case msg.Code == tea.KeyDown:
			m.editor.Move(
				editor.DirectionDown,
				editor.MoveByCharacter,
				msg.Mod.Contains(tea.ModShift),
			)

		case msg.Code == tea.KeySpace:
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
