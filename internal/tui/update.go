package tui

import (
	"fmt"
	"path/filepath"

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
		m.updateSearchInputWidth()

	case textSearchResultMsg:
		m.handleTextSearchResult(msg)
		m.updateViewport()

	case fileSearchResultMsg:
		m.handleFileSearchResult(msg)
		m.updateViewport()

	case tea.PasteMsg:
		if m.search.Mode != 0 {
			return m, m.handleSearchPaste(msg)
		}

		m.editor.InsertText(msg.String())
		m.clearMessage()
		m.updateViewport()

	case tea.MouseWheelMsg:
		if m.search.Mode != 0 {
			return m, nil
		}

		switch msg.Button {
		case tea.MouseWheelUp:
			m.viewport.ScrollUp(3)

		case tea.MouseWheelDown:
			m.viewport.ScrollDown(
				3,
				m.editor.LineCount(),
			)
		}

	case tea.KeyPressMsg:
		if m.search.Mode != 0 {
			cmd, quit := m.handleSearchKey(msg)

			m.updateViewport()

			if quit {
				return m, tea.Quit
			}

			return m, cmd
		}

		cmd, quit := m.handleKey(msg)

		m.updateViewport()

		if quit {
			return m, tea.Quit
		}

		return m, cmd
	}

	return m, nil
}

// handleKey handles editor keyboard input.
func (m *Model) handleKey(
	msg tea.KeyPressMsg,
) (tea.Cmd, bool) {
	shift := msg.Mod.Contains(tea.ModShift)
	ctrl := msg.Mod.Contains(tea.ModCtrl)

	switch {
	case ctrl && msg.Code == 'f':
		return m.openTextSearch(), false

	case ctrl && msg.Code == 'a':
		m.editor.SelectAll()
		m.clearMessage()

	case ctrl && msg.Code == 's':
		m.save()

	case ctrl && msg.Code == 'c':
		if text, ok := m.editor.SelectedText(); ok {
			if err := m.clipboard.Copy(text); err != nil {
				m.setMessage(
					MessageError,
					fmt.Sprintf(
						"copy failed: %v",
						err,
					),
				)
			} else {
				m.setMessage(
					MessageSuccess,
					"Copied",
				)
			}
		}

	case ctrl && msg.Code == 'x':
		if text, ok := m.editor.CutSelection(); ok {
			if err := m.clipboard.Copy(text); err != nil {
				m.setMessage(
					MessageError,
					fmt.Sprintf(
						"cut failed: %v",
						err,
					),
				)
			} else {
				m.setMessage(
					MessageSuccess,
					"Cut",
				)
			}
		}

	case ctrl && msg.Code == 'v':
		text, err := m.clipboard.Paste()
		if err != nil {
			m.setMessage(
				MessageError,
				fmt.Sprintf(
					"paste failed: %v",
					err,
				),
			)
			break
		}

		m.editor.InsertText(text)
		m.clearMessage()

	case ctrl && msg.Code == 'z' && shift:
		m.editor.Redo()
		m.clearMessage()

	case ctrl && msg.Code == 'z':
		m.editor.Undo()
		m.clearMessage()

	case ctrl && msg.Code == 'y':
		m.editor.Redo()
		m.clearMessage()

	case msg.Code == tea.KeyEscape:
		return nil, true

	case ctrl && msg.Code == tea.KeyLeft:
		m.editor.Move(
			editor.DirectionLeft,
			editor.MoveByWord,
			shift,
		)
		m.clearMessage()

	case ctrl && msg.Code == tea.KeyRight:
		m.editor.Move(
			editor.DirectionRight,
			editor.MoveByWord,
			shift,
		)
		m.clearMessage()

	case ctrl && msg.Code == 'h':
		m.editor.DeleteWordBackwards()
		m.clearMessage()

	case msg.Code == tea.KeyBackspace:
		m.editor.Backspace()
		m.clearMessage()

	case msg.Code == tea.KeyEnter:
		m.editor.Enter()
		m.clearMessage()

	case msg.Code == tea.KeyLeft:
		m.editor.Move(
			editor.DirectionLeft,
			editor.MoveByCharacter,
			shift,
		)
		m.clearMessage()

	case msg.Code == tea.KeyRight:
		m.editor.Move(
			editor.DirectionRight,
			editor.MoveByCharacter,
			shift,
		)
		m.clearMessage()

	case msg.Code == tea.KeyUp:
		m.editor.Move(
			editor.DirectionUp,
			editor.MoveByCharacter,
			shift,
		)
		m.clearMessage()

	case msg.Code == tea.KeyDown:
		m.editor.Move(
			editor.DirectionDown,
			editor.MoveByCharacter,
			shift,
		)
		m.clearMessage()

	case msg.Code == tea.KeySpace:
		m.editor.Insert(' ')
		m.clearMessage()

	default:
		if msg.Text != "" {
			m.editor.InsertText(msg.Text)
			m.clearMessage()
		}
	}

	return nil, false
}

func (m *Model) save() {
	if m.editor.FilePath() == "" {
		m.setMessage(
			MessageError,
			"cannot save: no file is associated with this document",
		)
		return
	}

	if !m.editor.Modified() {
		m.setMessage(
			MessageInfo,
			fmt.Sprintf(
				"%s is already saved",
				filepath.Base(
					m.editor.FilePath(),
				),
			),
		)
		return
	}

	if err := m.editor.Save(); err != nil {
		m.setMessage(
			MessageError,
			fmt.Sprintf(
				"save failed: %v",
				err,
			),
		)
		return
	}

	m.setMessage(
		MessageSuccess,
		fmt.Sprintf(
			"Saved %q",
			filepath.Base(
				m.editor.FilePath(),
			),
		),
	)
}

func (m *Model) clearMessage() {
	m.message = nil
}

func (m *Model) setMessage(
	level MessageLevel,
	text string,
) {
	m.message = &Message{
		Level: level,
		Text:  text,
	}
}

func (m *Model) updateViewport() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	lineNumberWidth := lineNumberWidth(
		m.editor.LineCount(),
	)

	searchHeight := m.searchHeight()

	messageHeight := 0

	if m.message != nil {
		messageHeight = 1
	}

	statusHeight := 1

	contentHeight := max(
		m.height-
			searchHeight-
			messageHeight-
			statusHeight,
		0,
	)

	m.viewport.SetSize(
		m.width,
		contentHeight,
	)

	cursor := m.editor.Cursor()

	m.viewport.EnsureCursorVisible(
		cursor,
		m.editor.Line(cursor.Line),
		m.editor.LineCount(),
		lineNumberWidth,
	)
}
