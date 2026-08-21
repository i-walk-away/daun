package main

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/i-walk-away/daun/internal/buffer"
	"github.com/i-walk-away/daun/internal/clipboard"
	"github.com/i-walk-away/daun/internal/editor"
	"github.com/i-walk-away/daun/internal/tui"
)

func main() {
	daunBuffer, err := loadBuffer(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	daunEditor := editor.NewEditor(daunBuffer)
	daunClipboard := clipboard.NewXClip()

	model := tui.NewModel(daunEditor, daunClipboard)

	_, err = tea.NewProgram(model).Run()

	closeErr := daunEditor.Close()

	if err != nil {
		log.Fatal(err)
	}

	if closeErr != nil {
		log.Fatal(closeErr)
	}
}

func loadBuffer(args []string) (*buffer.Buffer, error) {
	switch len(args) {
	case 0:
		return buffer.NewBuffer(), nil

	case 1:
		return openBuffer(args[0])

	default:
		return nil, fmt.Errorf("usage: daun [file]")
	}
}

func openBuffer(path string) (*buffer.Buffer, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %s", path)
		}

		return nil, fmt.Errorf("stat %q: %w", path, err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory: %s", path)
	}

	return buffer.NewBufferFromFile(path)
}
