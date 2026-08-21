package tui

import (
	"strings"

	"github.com/i-walk-away/daun/internal/editor"
)

// renderLineSearch renders a visible portion of a line and highlights every
// search match contained in that visible portion.
//
// The function performs all viewport calculations before applying Lip Gloss
// styles. Styled strings must never be passed to text slicing logic because
// they contain terminal escape sequences.
func renderLineSearch(
	text string,
	line int,
	query string,
	current *editor.Match,
	leftColumn int,
	contentWidth int,
) string {
	if contentWidth <= 0 {
		return ""
	}

	runes := []rune(text)

	if len(runes) == 0 {
		return ""
	}

	queryRunes := []rune(query)

	// Multiline queries are not highlighted by this line-level renderer.
	// The editor search still supports them; rendering them requires a
	// document-level highlight pass rather than independent line rendering.
	if len(queryRunes) == 0 ||
		strings.ContainsRune(query, '\n') {
		start, end := searchVisibleRuneRange(
			runes,
			leftColumn,
			contentWidth,
		)

		if start >= end {
			return ""
		}

		return string(runes[start:end])
	}

	start, end := searchVisibleRuneRange(
		runes,
		leftColumn,
		contentWidth,
	)

	if start >= end {
		return ""
	}

	var builder strings.Builder

	cursor := start

	for cursor < end {
		matchEnd := cursor + len(queryRunes)

		if matchEnd <= end &&
			searchRunesEqual(
				runes[cursor:matchEnd],
				queryRunes,
			) {

			style := searchMatchStyle

			if current != nil &&
				current.Start.Line == line &&
				current.Start.Column == cursor {

				style = searchCurrentMatchStyle
			}

			builder.WriteString(
				style.Render(
					string(runes[cursor:matchEnd]),
				),
			)

			cursor = matchEnd
			continue
		}

		builder.WriteRune(runes[cursor])
		cursor++
	}

	return builder.String()
}

func searchVisibleRuneRange(
	runes []rune,
	leftColumn int,
	contentWidth int,
) (int, int) {
	if len(runes) == 0 || contentWidth <= 0 {
		return 0, 0
	}

	leftColumn = max(leftColumn, 0)

	start := 0
	consumedWidth := 0

	for start < len(runes) {
		width := displayWidth(
			[]rune{runes[start]},
		)

		if consumedWidth+width > leftColumn {
			break
		}

		consumedWidth += width
		start++
	}

	end := start
	visibleWidth := 0

	for end < len(runes) {
		width := displayWidth(
			[]rune{runes[end]},
		)

		if visibleWidth+width > contentWidth {
			break
		}

		visibleWidth += width
		end++
	}

	return start, end
}

func searchRunesEqual(
	a []rune,
	b []rune,
) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
