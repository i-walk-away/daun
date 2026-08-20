package main

import (
	"log"

	tea "charm.land/bubbletea/v2"

	"github.com/i-walk-away/daun/internal/editor"
	"github.com/i-walk-away/daun/internal/tui"
)

func main() {
	buffer := editor.NewBuffer()
	buffer.SetLine(0, "Hello, Daun!")

	daunEditor := editor.NewEditor(buffer)
	model := tui.NewModel(daunEditor)

	if _, err := tea.NewProgram(model).Run(); err != nil {
		log.Fatal(err)
	}
}
