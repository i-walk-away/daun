package tui

import tea "charm.land/bubbletea/v2"

// Update handles incoming terminal events.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {

		// Keybinds are defined here
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
