package main

import (
	"log"

	tea "charm.land/bubbletea/v2"

	"github.com/i-walk-away/daun/internal/buffer"
	"github.com/i-walk-away/daun/internal/clipboard"
	"github.com/i-walk-away/daun/internal/editor"
	"github.com/i-walk-away/daun/internal/tui"
)

func main() {
	daunBuffer := buffer.NewBuffer()
	daunEditor := editor.NewEditor(daunBuffer)
	daunClipboard := clipboard.NewXClip()

	model := tui.NewModel(daunEditor, daunClipboard)

	if _, err := tea.NewProgram(model).Run(); err != nil {
		log.Fatal(err)
	}
}
