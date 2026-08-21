package tui

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/i-walk-away/daun/internal/buffer"
	"github.com/i-walk-away/daun/internal/editor"
)

func newTestModel(text string) Model {
	b := buffer.NewBuffer()
	e := editor.NewEditor(b)

	e.InsertText(text)

	return NewModel(e)
}

func TestCtrlASelectsAll(t *testing.T) {
	m := newTestModel("hello world")

	msg := tea.KeyPressMsg{
		Code: 'a',
		Mod:  tea.ModCtrl,
	}

	next, _ := m.Update(msg)
	m = next.(Model)

	start, end, ok := m.editor.Selection()

	if !ok {
		t.Fatal("expected selection")
	}

	if start != (editor.Position{}) {
		t.Fatalf("expected start at origin, got %+v", start)
	}

	if end != (editor.Position{Line: 0, Column: 11}) {
		t.Fatalf(
			"expected end %+v, got %+v",
			editor.Position{Line: 0, Column: 11},
			end,
		)
	}
}

func TestCtrlCProducesClipboardCommand(t *testing.T) {
	m := newTestModel("hello world")
	m.editor.SelectAll()

	msg := tea.KeyPressMsg{
		Code: 'c',
		Mod:  tea.ModCtrl,
	}

	next, cmd := m.Update(msg)

	if next == nil {
		t.Fatal("expected model")
	}

	if cmd == nil {
		t.Fatal("expected clipboard command")
	}
}

func TestCtrlXDeletesSelection(t *testing.T) {
	m := newTestModel("hello world")
	m.editor.SelectAll()

	msg := tea.KeyPressMsg{
		Code: 'x',
		Mod:  tea.ModCtrl,
	}

	next, cmd := m.Update(msg)
	m = next.(Model)

	if cmd == nil {
		t.Fatal("expected clipboard command")
	}

	if got := m.editor.Line(0); got != "" {
		t.Fatalf("expected empty line, got %q", got)
	}
}

func TestCtrlZUndoesEdit(t *testing.T) {
	m := newTestModel("hello")

	msg := tea.KeyPressMsg{
		Code: 'z',
		Mod:  tea.ModCtrl,
	}

	// Insert another character first.
	insert := tea.KeyPressMsg{
		Code: '!',
		Text: "!",
	}

	next, _ := m.Update(insert)
	m = next.(Model)

	if got := m.editor.Line(0); got != "hello!" {
		t.Fatalf("expected %q, got %q", "hello!", got)
	}

	next, _ = m.Update(msg)
	m = next.(Model)

	if got := m.editor.Line(0); got != "hello" {
		t.Fatalf("expected %q after undo, got %q", "hello", got)
	}
}

func TestCtrlYRedoesEdit(t *testing.T) {
	m := newTestModel("hello")

	insert := tea.KeyPressMsg{
		Code: '!',
		Text: "!",
	}

	next, _ := m.Update(insert)
	m = next.(Model)

	undo := tea.KeyPressMsg{
		Code: 'z',
		Mod:  tea.ModCtrl,
	}

	next, _ = m.Update(undo)
	m = next.(Model)

	redo := tea.KeyPressMsg{
		Code: 'y',
		Mod:  tea.ModCtrl,
	}

	next, _ = m.Update(redo)
	m = next.(Model)

	if got := m.editor.Line(0); got != "hello!" {
		t.Fatalf("expected %q, got %q", "hello!", got)
	}
}

func TestPasteMessageInsertsText(t *testing.T) {
	m := newTestModel("hello")

	msg := tea.PasteMsg{
		Content: " world",
	}

	next, _ := m.Update(msg)
	m = next.(Model)

	if got := m.editor.Line(0); got != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", got)
	}
}

func TestMouseWheelDoesNotMoveCursor(t *testing.T) {
	m := newTestModel("hello")

	before := m.editor.Cursor()

	msg := tea.MouseWheelMsg{
		Button: tea.MouseWheelDown,
	}

	next, _ := m.Update(msg)
	m = next.(Model)

	if got := m.editor.Cursor(); got != before {
		t.Fatalf(
			"mouse wheel moved cursor: expected %+v, got %+v",
			before,
			got,
		)
	}
}
