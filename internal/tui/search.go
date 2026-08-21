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

type SearchState struct {
	Mode  searchMode
	Input textinput.Model

	originalCursor editor.Position

	text  TextSearchState
	files FileSearchState
}

type TextSearchState struct {
	Current editor.Match

	CurrentIndex int
	Total        int
	Found        bool
	Wrapped      bool

	generation uint64
	cancel     context.CancelFunc
}

type FileSearchState struct {
	Results  []fileio.FileMatch
	Selected int
	Loading  bool

	generation uint64
	cancel     context.CancelFunc
}

type textSearchResultMsg struct {
	generation uint64
	query      string
	result     editor.SearchResult
	found      bool
}

type fileSearchResultMsg struct {
	generation uint64
	query      string
	results    []fileio.FileMatch
	err        error
}

func newTextSearchResult(
	generation uint64,
	query string,
	editor *editor.Editor,
	start editor.Position,
) tea.Cmd {
	return func() tea.Msg {
		result, found := editor.Search(
			query,
			start,
		)

		return textSearchResultMsg{
			generation: generation,
			query:      query,
			result:     result,
			found:      found,
		}
	}
}

func (m *Model) openTextSearch() tea.Cmd {
	m.cancelSearchOperations()

	m.search.Mode = searchModeText
	m.search.Input.Reset()
	m.search.Input.Prompt = "/ "
	m.search.Input.Placeholder = "Search text..."
	m.search.Input.SetVirtualCursor(false)

	m.search.originalCursor = m.editor.Cursor()

	if _, _, ok := m.editor.Selection(); ok {
		if selected, exists := m.editor.SelectedText(); exists {
			if len([]rune(selected)) > searchQueryLimit {
				selected = string(
					[]rune(selected)[:searchQueryLimit],
				)
			}

			m.search.Input.SetValue(selected)
		}
	}

	m.search.text = TextSearchState{}

	m.updateSearchInputWidth()
	m.clearMessage()

	return m.search.Input.Focus()
}

func (m *Model) openFileSearch() tea.Cmd {
	m.cancelSearchOperations()

	m.search.Mode = searchModeFiles
	m.search.Input.Reset()
	m.search.Input.Prompt = "> "
	m.search.Input.Placeholder = "Find files in home..."
	m.search.Input.SetVirtualCursor(false)

	m.search.files = FileSearchState{}

	m.updateSearchInputWidth()
	m.clearMessage()

	return nil
}

func (m *Model) closeSearch() {
	m.cancelSearchOperations()

	m.search.Mode = 0
	m.search.Input.Blur()
	m.search.Input.Reset()

	m.search.text = TextSearchState{}
	m.search.files = FileSearchState{}
}

func (m *Model) handleSearchKey(
	msg tea.KeyPressMsg,
) (tea.Cmd, bool) {
	ctrl := msg.Mod.Contains(tea.ModCtrl)
	shift := msg.Mod.Contains(tea.ModShift)

	switch {
	case ctrl && msg.Code == 'f':
		if m.search.Mode == searchModeText {
			return m.openFileSearch(), false
		}

		return m.openTextSearch(), false

	case msg.Code == tea.KeyEscape:
		m.closeSearch()
		return nil, false

	case m.search.Mode == searchModeText &&
		msg.Code == tea.KeyEnter:
		m.acceptTextSearch()
		return nil, false

	case m.search.Mode == searchModeText &&
		msg.Code == tea.KeyUp:
		m.navigateTextSearch(-1)
		return nil, false

	case m.search.Mode == searchModeText &&
		msg.Code == tea.KeyDown:
		m.navigateTextSearch(1)
		return nil, false

	case m.search.Mode == searchModeText &&
		ctrl &&
		msg.Code == 'p':
		m.navigateTextSearch(-1)
		return nil, false

	case m.search.Mode == searchModeText &&
		ctrl &&
		msg.Code == 'n':
		m.navigateTextSearch(1)
		return nil, false

	case m.search.Mode == searchModeFiles &&
		msg.Code == tea.KeyUp:
		m.moveFileSelection(-1)
		return nil, false

	case m.search.Mode == searchModeFiles &&
		msg.Code == tea.KeyDown:
		m.moveFileSelection(1)
		return nil, false

	case m.search.Mode == searchModeFiles &&
		msg.Code == tea.KeyEnter:
		m.acceptFileSearch()
		return nil, false

	case m.search.Mode == searchModeText &&
		shift &&
		msg.Code == tea.KeyEnter:
		m.navigateTextSearch(-1)
		return nil, false
	}

	input, cmd := m.search.Input.Update(msg)
	m.search.Input = input

	if m.search.Mode == searchModeText {
		return tea.Batch(
			cmd,
			m.startTextSearch(),
		), false
	}

	return tea.Batch(
		cmd,
		m.startFileSearch(),
	), false
}

func (m *Model) handleSearchPaste(
	msg tea.PasteMsg,
) tea.Cmd {
	input, cmd := m.search.Input.Update(msg)

	m.search.Input = input

	if m.search.Mode == searchModeText {
		return tea.Batch(
			cmd,
			m.startTextSearch(),
		)
	}

	return tea.Batch(
		cmd,
		m.startFileSearch(),
	)
}

func (m *Model) startTextSearch() tea.Cmd {
	m.cancelTextSearch()

	query := m.search.Input.Value()

	if query == "" {
		m.search.text = TextSearchState{}
		return nil
	}

	m.search.text.generation++
	generation := m.search.text.generation

	start := m.search.originalCursor
	editorModel := m.editor

	return newTextSearchResult(
		generation,
		query,
		editorModel,
		start,
	)
}

func (m *Model) handleTextSearchResult(
	msg textSearchResultMsg,
) {
	if m.search.Mode != searchModeText {
		return
	}

	if msg.generation != m.search.text.generation {
		return
	}

	if msg.query != m.search.Input.Value() {
		return
	}

	if !msg.found {
		m.search.text = TextSearchState{
			generation: msg.generation,
		}

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

	var position editor.Position
	var found bool

	if delta > 0 {
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

		if !found {
			return
		}

		next, ok := m.editor.Search(
			query,
			position,
		)
		if !ok {
			return
		}

		m.search.text.Current = next.Match

		if m.search.text.Total > 0 {
			m.search.text.CurrentIndex++
			if m.search.text.CurrentIndex >
				m.search.text.Total {
				m.search.text.CurrentIndex = 1
			}
		}

		return
	}

	position, found = m.editor.FindPrevious(
		query,
		m.search.text.Current.Start,
	)

	if !found {
		// Wrap to the last match.
		position, found = m.findLastMatch(query)
	}

	if !found {
		return
	}

	previous, ok := m.editor.Search(
		query,
		position,
	)

	if !ok {
		return
	}

	m.search.text.Current = previous.Match

	if m.search.text.Total > 0 {
		m.search.text.CurrentIndex--
		if m.search.text.CurrentIndex <= 0 {
			m.search.text.CurrentIndex =
				m.search.text.Total
		}
	}
}

func (m *Model) findLastMatch(
	query string,
) (editor.Position, bool) {
	lastLine := m.editor.LineCount() - 1

	cursor := editor.Position{
		Line:   lastLine,
		Column: len([]rune(m.editor.Line(lastLine))),
	}

	return m.editor.FindPrevious(
		query,
		cursor,
	)
}

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
	rootFinder := m.fileFinder

	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	m.search.files.cancel = cancel
	m.search.files.Loading = true

	return func() tea.Msg {
		timer := time.NewTimer(
			searchInputDebounce,
		)

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

		results, err := rootFinder.Search(
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

	if msg.query !=
		strings.TrimSpace(
			m.search.Input.Value(),
		) {
		return
	}

	m.search.files.cancel = nil
	m.search.files.Loading = false

	if msg.err != nil {
		if msg.err == context.Canceled {
			return
		}

		m.setMessage(
			MessageError,
			"file search failed: "+msg.err.Error(),
		)

		m.search.files.Results = nil
		m.search.files.Selected = 0

		return
	}

	m.search.files.Results = msg.results

	if m.search.files.Selected >=
		len(msg.results) {
		m.search.files.Selected = max(
			len(msg.results)-1,
			0,
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

func (m *Model) cancelTextSearch() {
	if m.search.text.cancel != nil {
		m.search.text.cancel()
		m.search.text.cancel = nil
	}
}

func (m *Model) cancelFileSearch() {
	if m.search.files.cancel != nil {
		m.search.files.cancel()
		m.search.files.cancel = nil
	}

	m.search.files.Loading = false
}

func (m *Model) cancelSearchOperations() {
	m.cancelTextSearch()
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
		// Separator + input.
		return 2

	case searchModeFiles:
		height := 2

		if len(m.search.files.Results) == 0 &&
			!m.search.files.Loading &&
			m.search.Input.Value() != "" {
			height++
		}

		return height + min(
			len(m.search.files.Results),
			fileVisibleLimit,
		)

	default:
		return 0
	}
}

func (m Model) renderSearchPanel() string {
	if m.search.Mode == 0 {
		return ""
	}

	var b strings.Builder

	b.WriteString(
		searchSeparatorStyle.Render(
			strings.Repeat(
				"─",
				max(m.width, 1),
			),
		),
	)

	b.WriteByte('\n')

	inputView := m.search.Input.View()

	if m.search.Mode == searchModeText {
		counter := ""

		switch {
		case !m.search.text.Found:
			if m.search.Input.Value() != "" {
				counter = "No matches"
			}

		default:
			counter = formatSearchCounter(
				m.search.text.CurrentIndex,
				m.search.text.Total,
			)
		}

		inputWidth := len([]rune(inputView))
		counterWidth := len([]rune(counter))

		padding := max(
			m.width-inputWidth-counterWidth,
			1,
		)

		b.WriteString(inputView)
		b.WriteString(strings.Repeat(" ", padding))
		b.WriteString(counter)

		return b.String()
	}

	b.WriteString(inputView)

	if m.search.files.Loading {
		b.WriteString("  " + searchInfoStyle.Render("…"))
	}

	for i, match := range m.search.files.Results {
		if i >= fileVisibleLimit {
			break
		}

		b.WriteByte('\n')

		style := searchResultStyle

		if i == m.search.files.Selected {
			style = searchResultSelectedStyle
		}

		b.WriteString(
			style.Width(
				max(m.width-2, 1),
			).Render(
				"  " + match.DisplayPath,
			),
		)
	}

	if len(m.search.files.Results) == 0 &&
		!m.search.files.Loading &&
		m.search.Input.Value() != "" {

		b.WriteByte('\n')
		b.WriteString(
			searchInfoStyle.Render(
				"  No files found",
			),
		)
	}

	return b.String()
}

func formatSearchCounter(
	current int,
	total int,
) string {
	return fmt.Sprintf(
		"%d / %d",
		current,
		total,
	)
}
