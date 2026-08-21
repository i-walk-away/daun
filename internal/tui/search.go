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

	generation uint64
	cancel     context.CancelFunc
}

// FileSearchState contains file-picker state.
type FileSearchState struct {
	Results []fileio.FileMatch

	Selected     int
	ScrollOffset int

	Loading bool

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
	editorModel *editor.Editor,
	start editor.Position,
) tea.Cmd {
	return func() tea.Msg {
		result, found := editorModel.Search(
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
	m.search.text = TextSearchState{}

	if _, _, ok := m.editor.Selection(); ok {
		if selected, exists := m.editor.SelectedText(); exists {
			runes := []rune(selected)

			if len(runes) > searchQueryLimit {
				selected = string(runes[:searchQueryLimit])
			}

			m.search.Input.SetValue(selected)
		}
	}

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

		if shift {
			m.navigateTextSearch(-1)
		} else {
			m.acceptTextSearch()
		}

		return nil, false

	case m.search.Mode == searchModeText &&
		msg.Code == tea.KeyUp:

		m.navigateTextSearch(-1)
		return nil, false

	case m.search.Mode == searchModeText &&
		msg.Code == tea.KeyDown:

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
	}

	input, cmd := m.search.Input.Update(msg)
	m.search.Input = input

	switch m.search.Mode {
	case searchModeText:
		return tea.Batch(
			cmd,
			m.startTextSearch(),
		), false

	case searchModeFiles:
		return tea.Batch(
			cmd,
			m.startFileSearch(),
		), false
	}

	return cmd, false
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

	return newTextSearchResult(
		generation,
		query,
		m.editor,
		m.search.originalCursor,
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

	var (
		position editor.Position
		found    bool
	)

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
	} else {
		position, found = m.editor.FindPrevious(
			query,
			m.search.text.Current.Start,
		)

		if !found {
			position, found = m.findLastMatch(query)
		}

		if !found {
			return
		}
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

	m.scrollCursorIntoView(result.Match.Start)
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

	maxVisible := fileVisibleLimit

	if selected >= m.search.files.ScrollOffset+maxVisible {
		m.search.files.ScrollOffset =
			selected - maxVisible + 1
	}

	maxOffset := max(
		len(m.search.files.Results)-maxVisible,
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

			if i == m.search.files.Selected {
				style = searchResultSelectedStyle
			}

			prefix := "  "

			if i == m.search.files.Selected {
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
