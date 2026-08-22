package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/i-walk-away/daun/internal/editor"
	"github.com/i-walk-away/daun/internal/fileio"
)

const (
	fileSearchLimit     = 100
	fileVisibleLimit    = 8
	searchQueryLimit    = 4096
	searchInputDebounce = 100 * time.Millisecond
)

type searchMode uint8

const (
	searchModeText searchMode = iota + 1
	searchModeFiles
)

// SearchState contains the state of the search overlay.
type SearchState struct {
	Mode  searchMode
	Input textinput.Model

	originalCursor editor.Position

	text  TextSearchState
	files FileSearchState
}

// TextSearchState contains text-search state.
type TextSearchState struct {
	Current editor.Match

	CurrentIndex int
	Total        int
	Found        bool
	Wrapped      bool
}

// FileSearchState contains file-search state.
type FileSearchState struct {
	Results []fileio.FileMatch

	Selected     int
	ScrollOffset int

	Loading bool

	generation uint64
	cancel     context.CancelFunc
}

type textSearchResultMsg struct {
	query  string
	result editor.SearchResult
	found  bool
}

type fileSearchResultMsg struct {
	generation uint64
	query      string
	results    []fileio.FileMatch
	err        error
}

// openTextSearch opens text search and pre-fills it with the current
// selection when one exists.
func (m *Model) openTextSearch() tea.Cmd {
	m.beginSearch(
		searchModeText,
		"/ ",
		"Search text...",
	)

	m.search.originalCursor = m.editor.Cursor()
	m.search.text = TextSearchState{}

	if _, _, ok := m.editor.Selection(); ok {
		selected, exists := m.editor.SelectedText()
		if exists {
			runes := []rune(selected)

			if len(runes) > searchQueryLimit {
				selected = string(runes[:searchQueryLimit])
			}

			m.search.Input.SetValue(selected)
		}
	}

	return m.search.Input.Focus()
}

// openFileSearch opens the file picker.
func (m *Model) openFileSearch() tea.Cmd {
	m.beginSearch(
		searchModeFiles,
		"> ",
		"Find files in home...",
	)

	m.search.files = FileSearchState{}

	return nil
}

// beginSearch resets common search UI state before entering a search mode.
func (m *Model) beginSearch(
	mode searchMode,
	prompt string,
	placeholder string,
) {
	m.cancelSearchOperations()

	m.search.Mode = mode
	m.search.Input.Reset()
	m.search.Input.Prompt = prompt
	m.search.Input.Placeholder = placeholder
	m.search.Input.SetVirtualCursor(false)

	m.updateSearchInputWidth()
	m.clearMessage()
}

// closeSearch closes the search overlay and cancels any running file search.
func (m *Model) closeSearch() {
	m.cancelSearchOperations()

	m.search.Mode = 0
	m.search.Input.Blur()
	m.search.Input.Reset()

	m.search.text = TextSearchState{}
	m.search.files = FileSearchState{}
}

// handleSearchKey dispatches keyboard input according to the active search
// mode.
func (m *Model) handleSearchKey(
	msg tea.KeyPressMsg,
) (tea.Cmd, bool) {
	ctrl := msg.Mod.Contains(tea.ModCtrl)

	switch {
	case ctrl && msg.Code == 'f':
		if m.search.Mode == searchModeText {
			return m.openFileSearch(), false
		}

		return m.openTextSearch(), false

	case msg.Code == tea.KeyEscape:
		m.closeSearch()
		return nil, false
	}

	switch m.search.Mode {
	case searchModeText:
		return m.handleTextSearchKey(msg)

	case searchModeFiles:
		return m.handleFileSearchKey(msg)
	}

	return nil, false
}

func (m *Model) handleTextSearchKey(
	msg tea.KeyPressMsg,
) (tea.Cmd, bool) {
	switch msg.Code {
	case tea.KeyEnter:
		if msg.Mod.Contains(tea.ModShift) {
			m.navigateTextSearch(-1)
		} else {
			m.acceptTextSearch()
		}

		return nil, false

	case tea.KeyUp:
		m.navigateTextSearch(-1)
		return nil, false

	case tea.KeyDown:
		m.navigateTextSearch(1)
		return nil, false
	}

	input, cmd := m.search.Input.Update(msg)
	m.search.Input = input

	return tea.Batch(
		cmd,
		m.startTextSearch(),
	), false
}

func (m *Model) handleFileSearchKey(
	msg tea.KeyPressMsg,
) (tea.Cmd, bool) {
	switch msg.Code {
	case tea.KeyUp:
		m.moveFileSelection(-1)
		return nil, false

	case tea.KeyDown:
		m.moveFileSelection(1)
		return nil, false

	case tea.KeyEnter:
		m.acceptFileSearch()
		return nil, false
	}

	input, cmd := m.search.Input.Update(msg)
	m.search.Input = input

	return tea.Batch(
		cmd,
		m.startFileSearch(),
	), false
}

// handleSearchPaste inserts terminal paste into the active search input.
func (m *Model) handleSearchPaste(
	msg tea.PasteMsg,
) tea.Cmd {
	input, cmd := m.search.Input.Update(msg)
	m.search.Input = input

	switch m.search.Mode {
	case searchModeText:
		return tea.Batch(
			cmd,
			m.startTextSearch(),
		)

	case searchModeFiles:
		return tea.Batch(
			cmd,
			m.startFileSearch(),
		)
	}

	return cmd
}

// startTextSearch starts a text-search operation.
//
// Text search is intentionally synchronous with respect to command creation.
// Bubble Tea executes the returned command independently, while the editor
// remains unchanged until the result is received.
func (m *Model) startTextSearch() tea.Cmd {
	query := m.search.Input.Value()

	if query == "" {
		m.search.text = TextSearchState{}
		return nil
	}

	start := m.search.originalCursor
	editorModel := m.editor

	return func() tea.Msg {
		result, found := editorModel.Search(
			query,
			start,
		)

		return textSearchResultMsg{
			query:  query,
			result: result,
			found:  found,
		}
	}
}

func (m *Model) handleTextSearchResult(
	msg textSearchResultMsg,
) {
	if m.search.Mode != searchModeText {
		return
	}

	if msg.query != m.search.Input.Value() {
		return
	}

	if !msg.found {
		m.search.text = TextSearchState{}
		return
	}

	m.search.text.Current = msg.result.Match
	m.search.text.CurrentIndex = msg.result.Current
	m.search.text.Total = msg.result.Total
	m.search.text.Found = true
	m.search.text.Wrapped = msg.result.Wrapped
}

func (m *Model) navigateTextSearch(delta int) {
	if !m.search.text.Found {
		return
	}

	query := m.search.Input.Value()

	var (
		position editor.Position
		found    bool
	)

	switch {
	case delta > 0:
		position, found = m.editor.Find(
			query,
			m.search.text.Current.End,
		)

		if !found {
			position, found = m.editor.Find(
				query,
				editor.Position{},
			)
		}

	default:
		position, found = m.editor.FindPrevious(
			query,
			m.search.text.Current.Start,
		)

		if !found {
			position, found = m.findLastMatch(query)
		}
	}

	if !found {
		return
	}

	result, found := m.editor.Search(
		query,
		position,
	)
	if !found {
		return
	}

	m.search.text.Current = result.Match
	m.search.text.CurrentIndex = result.Current
	m.search.text.Total = result.Total
	m.search.text.Wrapped = result.Wrapped

	m.scrollCursorIntoView(
		result.Match.Start,
	)
}

func (m *Model) findLastMatch(
	query string,
) (editor.Position, bool) {
	lastLine := m.editor.LineCount() - 1

	lastPosition := editor.Position{
		Line: lastLine,
		Column: len(
			[]rune(m.editor.Line(lastLine)),
		),
	}

	return m.editor.FindPrevious(
		query,
		lastPosition,
	)
}

// acceptTextSearch moves the editor cursor to the end of the current match
// and closes the search overlay.
func (m *Model) acceptTextSearch() {
	if !m.search.text.Found {
		return
	}

	m.editor.SetCursor(
		m.search.text.Current.End,
	)

	m.closeSearch()
	m.updateViewport()
}

func (m *Model) startFileSearch() tea.Cmd {
	m.cancelFileSearch()

	query := strings.TrimSpace(
		m.search.Input.Value(),
	)

	if query == "" {
		m.search.files = FileSearchState{}
		return nil
	}

	m.search.files.generation++

	generation := m.search.files.generation

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	m.search.files.cancel = cancel
	m.search.files.Loading = true

	return func() tea.Msg {
		timer := time.NewTimer(searchInputDebounce)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return fileSearchResultMsg{
				generation: generation,
				query:      query,
				err:        ctx.Err(),
			}

		case <-timer.C:
		}

		results, err := m.fileFinder.Search(
			ctx,
			query,
			fileSearchLimit,
		)

		return fileSearchResultMsg{
			generation: generation,
			query:      query,
			results:    results,
			err:        err,
		}
	}
}

func (m *Model) handleFileSearchResult(
	msg fileSearchResultMsg,
) {
	if m.search.Mode != searchModeFiles {
		return
	}

	if msg.generation !=
		m.search.files.generation {
		return
	}

	if msg.query != strings.TrimSpace(
		m.search.Input.Value(),
	) {
		return
	}

	m.search.files.cancel = nil
	m.search.files.Loading = false
	m.search.files.Results = msg.results

	if len(msg.results) == 0 {
		m.search.files.Selected = 0
		m.search.files.ScrollOffset = 0
	} else {
		m.search.files.Selected = min(
			m.search.files.Selected,
			len(msg.results)-1,
		)

		m.ensureFileSelectionVisible()
	}

	if msg.err != nil && len(msg.results) == 0 {
		m.setMessage(
			MessageError,
			"file search failed: "+msg.err.Error(),
		)
	}
}

func (m *Model) moveFileSelection(delta int) {
	count := len(m.search.files.Results)

	if count == 0 {
		return
	}

	m.search.files.Selected += delta

	if m.search.files.Selected < 0 {
		m.search.files.Selected = count - 1
	}

	if m.search.files.Selected >= count {
		m.search.files.Selected = 0
	}

	m.ensureFileSelectionVisible()
}

func (m *Model) ensureFileSelectionVisible() {
	selected := m.search.files.Selected

	if selected < m.search.files.ScrollOffset {
		m.search.files.ScrollOffset = selected
	}

	if selected >=
		m.search.files.ScrollOffset+fileVisibleLimit {

		m.search.files.ScrollOffset =
			selected - fileVisibleLimit + 1
	}

	maxOffset := max(
		len(m.search.files.Results)-fileVisibleLimit,
		0,
	)

	m.search.files.ScrollOffset = min(
		m.search.files.ScrollOffset,
		maxOffset,
	)
}

func (m *Model) acceptFileSearch() {
	if len(m.search.files.Results) == 0 {
		if m.search.files.Loading {
			m.setMessage(
				MessageInfo,
				"Still searching files...",
			)
		}

		return
	}

	match := m.search.files.Results[m.search.files.Selected]

	if err := m.editor.OpenFile(match.Path); err != nil {
		m.setMessage(
			MessageError,
			err.Error(),
		)

		return
	}

	m.closeSearch()

	m.setMessage(
		MessageSuccess,
		"Opened "+match.DisplayPath,
	)

	m.updateViewport()
}

func (m *Model) scrollCursorIntoView(
	position editor.Position,
) {
	lineNumberWidth := lineNumberWidth(
		m.editor.LineCount(),
	)

	m.viewport.EnsureCursorVisible(
		position,
		m.editor.Line(position.Line),
		m.editor.LineCount(),
		lineNumberWidth,
	)
}

func (m *Model) cancelFileSearch() {
	if m.search.files.cancel != nil {
		m.search.files.cancel()
		m.search.files.cancel = nil
	}

	m.search.files.Loading = false
}

func (m *Model) cancelSearchOperations() {
	m.cancelFileSearch()
}

func (m *Model) updateSearchInputWidth() {
	if m.width <= 0 {
		return
	}

	m.search.Input.SetWidth(
		max(m.width-16, 20),
	)
}

func (m Model) searchHeight() int {
	switch m.search.Mode {
	case searchModeText:
		return 2

	case searchModeFiles:
		height := 2

		height += min(
			len(m.search.files.Results),
			fileVisibleLimit,
		)

		if len(m.search.files.Results) == 0 &&
			!m.search.files.Loading &&
			m.search.Input.Value() != "" {
			height++
		}

		return height

	default:
		return 0
	}
}

func (m Model) renderSearchPanel() string {
	if m.search.Mode == 0 {
		return ""
	}

	var builder strings.Builder

	builder.WriteString(
		searchSeparatorStyle.Render(
			strings.Repeat(
				"─",
				max(m.width, 1),
			),
		),
	)

	builder.WriteByte('\n')

	inputView := m.search.Input.View()

	switch m.search.Mode {
	case searchModeText:
		counter := ""

		switch {
		case m.search.text.Found:
			counter = fmt.Sprintf(
				"%d / %d",
				m.search.text.CurrentIndex,
				m.search.text.Total,
			)

		case m.search.Input.Value() != "":
			counter = "No matches"
		}

		inputWidth := len([]rune(inputView))
		counterWidth := len([]rune(counter))

		padding := max(
			m.width-inputWidth-counterWidth,
			1,
		)

		builder.WriteString(inputView)
		builder.WriteString(
			strings.Repeat(" ", padding),
		)
		builder.WriteString(counter)

	case searchModeFiles:
		builder.WriteString(inputView)

		if m.search.files.Loading {
			builder.WriteString(
				"  " + searchInfoStyle.Render("…"),
			)
		}

		start := m.search.files.ScrollOffset
		end := min(
			start+fileVisibleLimit,
			len(m.search.files.Results),
		)

		for i := start; i < end; i++ {
			match := m.search.files.Results[i]

			builder.WriteByte('\n')

			style := searchResultStyle
			prefix := "  "

			if i == m.search.files.Selected {
				style = searchResultSelectedStyle
				prefix = "› "
			}

			builder.WriteString(
				style.
					Width(max(m.width-2, 1)).
					Render(
						prefix + match.DisplayPath,
					),
			)
		}

		if len(m.search.files.Results) == 0 &&
			!m.search.files.Loading &&
			m.search.Input.Value() != "" {

			builder.WriteByte('\n')
			builder.WriteString(
				searchInfoStyle.Render(
					"  No files found",
				),
			)
		}
	}

	return builder.String()
}
