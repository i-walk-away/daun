package tui

import (
	tea "charm.land/bubbletea/v2"
)

// Update handles incoming terminal events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		// Keybinds are defined here
		case msg.Code == tea.KeyEscape:
			return m, tea.Quit

		case msg.Mod.Contains(tea.ModCtrl) && msg.Code == tea.KeyLeft:
			m.editor.MoveWordLeft()

		case msg.Mod.Contains(tea.ModCtrl) && msg.Code == tea.KeyRight:
			m.editor.MoveWordRight()

		case msg.Mod.Contains(tea.ModCtrl) && msg.Code == 'h':
			m.editor.DeleteWordBackwards()

		case msg.Code == tea.KeyBackspace:
			m.editor.Backspace()

		case msg.Code == tea.KeyEnter:
			m.editor.Enter()

		case msg.Code == tea.KeyLeft:
			m.editor.MoveLeft()

		case msg.Code == tea.KeyRight:
			m.editor.MoveRight()

		case msg.Code == tea.KeyUp:
			m.editor.MoveUp()

		case msg.Code == tea.KeyDown:
			m.editor.MoveDown()

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
