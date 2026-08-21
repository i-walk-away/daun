package main

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/i-walk-away/daun/internal/buffer"
	"github.com/i-walk-away/daun/internal/clipboard"
	"github.com/i-walk-away/daun/internal/editor"
	"github.com/i-walk-away/daun/internal/fileio"
	"github.com/i-walk-away/daun/internal/tui"
)

func main() {
	daunEditor, err := loadEditor(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := daunEditor.Close(); err != nil {
			log.Printf("close editor: %v", err)
		}
	}()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	daunClipboard := clipboard.NewXClip()
	fileFinder := fileio.NewFileFinder(homeDir)

	model := tui.NewModel(
		daunEditor,
		daunClipboard,
		fileFinder,
	)

	if _, err := tea.NewProgram(model).Run(); err != nil {
		log.Fatal(err)
	}
}

func loadEditor(args []string) (*editor.Editor, error) {
	switch len(args) {
	case 0:
		return editor.NewEditor(
			buffer.NewBuffer(),
		), nil

	case 1:
		return loadFileEditor(args[0])

	default:
		return nil, fmt.Errorf("usage: daun [file]")
	}
}

func loadFileEditor(path string) (*editor.Editor, error) {
	buf, err := fileio.OpenOrCreate(path)
	if err != nil {
		return nil, err
	}

	return editor.NewFileEditor(
		buf,
		path,
	), nil
}
