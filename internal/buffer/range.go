package buffer

import "strings"

// TextRange returns the text contained in the specified range.
//
// The range is half-open: start is inclusive and end is exclusive.
// Newline characters between lines are included in the returned text.
func (b *Buffer) TextRange(
	startLine, startColumn,
	endLine, endColumn int,
) string {
	if startLine == endLine {
		line := []rune(b.Line(startLine))

		startColumn = min(max(startColumn, 0), len(line))
		endColumn = min(max(endColumn, startColumn), len(line))

		return string(line[startColumn:endColumn])
	}

	first := []rune(b.Line(startLine))
	last := []rune(b.Line(endLine))

	startColumn = min(max(startColumn, 0), len(first))
	endColumn = min(max(endColumn, 0), len(last))

	var result strings.Builder

	result.WriteString(string(first[startColumn:]))
	result.WriteByte('\n')

	for line := startLine + 1; line < endLine; line++ {
		result.WriteString(b.Line(line))
		result.WriteByte('\n')
	}

	result.WriteString(string(last[:endColumn]))

	return result.String()
}

// ReplaceRange replaces the specified text range with replacement.
//
// The range is half-open: start is inclusive and end is exclusive.
// Replacement may contain newline characters.
func (b *Buffer) ReplaceRange(
	startLine, startColumn,
	endLine, endColumn int,
	replacement string,
) {
	first := []rune(b.Line(startLine))
	last := []rune(b.Line(endLine))

	startColumn = min(max(startColumn, 0), len(first))
	endColumn = min(max(endColumn, 0), len(last))

	prefix := string(first[:startColumn])
	suffix := string(last[endColumn:])

	parts := strings.Split(replacement, "\n")

	replacementLines := make([]Line, len(parts))

	if len(parts) == 1 {
		replacementLines[0] = newLine(prefix + parts[0] + suffix)
	} else {
		replacementLines[0] = newLine(prefix + parts[0])

		for i := 1; i < len(parts)-1; i++ {
			replacementLines[i] = newLine(parts[i])
		}

		replacementLines[len(parts)-1] = newLine(
			parts[len(parts)-1] + suffix,
		)
	}

	removedLines := endLine - startLine + 1
	newLines := make([]Line, 0, len(b.lines)-removedLines+len(replacementLines))

	newLines = append(newLines, b.lines[:startLine]...)
	newLines = append(newLines, replacementLines...)
	newLines = append(newLines, b.lines[endLine+1:]...)

	b.lines = newLines
}
