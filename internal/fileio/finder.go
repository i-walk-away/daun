package fileio

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// FileMatch represents a file found by a filesystem search.
type FileMatch struct {
	Path        string
	DisplayPath string

	score int
	depth int
}

// FileFinder searches regular files below a filesystem root.
type FileFinder struct {
	root   string
	rgPath string
}

// NewFileFinder creates a file finder rooted at root.
//
// Ripgrep is used only to enumerate filesystem paths. Matching and ranking
// are performed by Daun.
func NewFileFinder(root string) *FileFinder {
	rgPath, _ := exec.LookPath("rg")

	return &FileFinder{
		root:   filepath.Clean(root),
		rgPath: rgPath,
	}
}

// Search finds files whose relative path contains query.
//
// Matching is case-insensitive. Hidden and ignored files are included.
// Symbolic links are not followed.
//
// Permission errors encountered while traversing the filesystem do not
// fail the complete search. At most limit results are returned. A
// non-positive limit returns all matches.
//
// Search stops when ctx is cancelled.
func (f *FileFinder) Search(
	ctx context.Context,
	query string,
	limit int,
) ([]FileMatch, error) {
	query = strings.TrimSpace(query)

	if query == "" {
		return nil, nil
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if f.rgPath == "" {
		return nil, fmt.Errorf(
			"ripgrep is required for file search: install 'rg'",
		)
	}

	info, err := os.Stat(f.root)
	if err != nil {
		return nil, fmt.Errorf(
			"stat search root: %w",
			err,
		)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf(
			"search root is not a directory: %s",
			f.root,
		)
	}

	query = strings.ToLower(
		filepath.ToSlash(query),
	)

	cmd := exec.CommandContext(
		ctx,
		f.rgPath,
		"--no-config",
		"--files",
		"--hidden",
		"--no-ignore",
		"--color",
		"never",
		"--",
		f.root,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf(
			"create ripgrep stdout pipe: %w",
			err,
		)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf(
			"create ripgrep stderr pipe: %w",
			err,
		)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf(
			"start ripgrep: %w",
			err,
		)
	}

	results := make([]FileMatch, 0, max(limit, 16))

	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(
		make([]byte, 64*1024),
		1024*1024,
	)

	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			_ = cmd.Wait()
			return nil, err
		}

		absolutePath := filepath.Clean(
			scanner.Text(),
		)

		relativePath, err := filepath.Rel(
			f.root,
			absolutePath,
		)
		if err != nil {
			continue
		}

		if relativePath == "." {
			continue
		}

		displayPath := filepath.ToSlash(
			relativePath,
		)

		lowerPath := strings.ToLower(
			displayPath,
		)

		base := strings.ToLower(
			filepath.Base(displayPath),
		)

		score, matched := rankPath(
			query,
			lowerPath,
			base,
		)

		if !matched {
			continue
		}

		results = append(
			results,
			FileMatch{
				Path:        absolutePath,
				DisplayPath: displayPath,
				score:       score,
				depth:       pathDepth(displayPath),
			},
		)
	}

	if err := scanner.Err(); err != nil {
		_ = cmd.Wait()

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		return nil, fmt.Errorf(
			"read ripgrep output: %w",
			err,
		)
	}

	stderrBytes, stderrErr := io.ReadAll(stderr)
	waitErr := cmd.Wait()

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	stderrText := strings.TrimSpace(
		string(stderrBytes),
	)

	if stderrErr != nil && stderrText == "" {
		stderrText = stderrErr.Error()
	}

	if waitErr != nil {
		if isPermissionOnly(stderrText) {
			// Ripgrep found at least some files but encountered paths that
			// the current user cannot inspect. This is expected for a
			// home-directory-wide search and should not invalidate results.
		} else {
			if stderrText == "" {
				stderrText = waitErr.Error()
			}

			return nil, fmt.Errorf(
				"ripgrep failed: %s",
				stderrText,
			)
		}
	}

	sort.Slice(
		results,
		func(i, j int) bool {
			return better(
				results[i],
				results[j],
			)
		},
	)

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

func isPermissionOnly(stderr string) bool {
	if stderr == "" {
		return false
	}

	lines := strings.Split(
		stderr,
		"\n",
	)

	foundPermissionError := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		lower := strings.ToLower(line)

		if strings.Contains(
			lower,
			"permission denied",
		) {
			foundPermissionError = true
			continue
		}

		if strings.Contains(
			lower,
			"access is denied",
		) {
			foundPermissionError = true
			continue
		}

		// Ripgrep may emit "os error 13" for EACCES.
		if strings.Contains(
			lower,
			"os error 13",
		) {
			foundPermissionError = true
			continue
		}

		// Any unrelated stderr output means this is not a purely
		// permission-related failure.
		return false
	}

	return foundPermissionError
}

func rankPath(
	query string,
	path string,
	base string,
) (score int, matched bool) {
	switch {
	case base == query:
		return 10000, true

	case strings.TrimSuffix(
		base,
		filepath.Ext(base),
	) == query:
		return 9500, true

	case strings.HasPrefix(base, query):
		return 8500 - len(base), true

	case strings.Contains(base, query):
		return 7000 - len(base), true

	case strings.HasPrefix(path, query):
		return 5000 - len(path), true

	case strings.Contains(path, query):
		return 3000 - len(path), true
	}

	return 0, false
}

func pathDepth(path string) int {
	return strings.Count(path, "/")
}

func better(a, b FileMatch) bool {
	if a.score != b.score {
		return a.score > b.score
	}

	if a.depth != b.depth {
		return a.depth < b.depth
	}

	return a.DisplayPath < b.DisplayPath
}
