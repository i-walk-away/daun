package buffer

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// NewBufferFromFile opens a file using a memory-mapped original source.
//
// The file contents are not copied into the Go heap. The operating system
// manages the mapped pages and loads them into memory as they are accessed.
//
// The returned buffer owns the mapping and must be closed when it is no longer
// needed.
func NewBufferFromFile(path string) (*Buffer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	info, err := file.Stat()
	if err != nil {
		_ = file.Close()

		return nil, fmt.Errorf("stat file: %w", err)
	}

	if !info.Mode().IsRegular() {
		_ = file.Close()

		return nil, fmt.Errorf("path is not a regular file: %s", path)
	}

	size := info.Size()

	if size > int64(int(^uint(0)>>1)) {
		_ = file.Close()

		return nil, fmt.Errorf("file is too large: %d bytes", size)
	}

	if size == 0 {
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("close empty file: %w", err)
		}

		return NewBuffer(), nil
	}

	mapped, err := unix.Mmap(
		int(file.Fd()),
		0,
		int(size),
		unix.PROT_READ,
		unix.MAP_PRIVATE,
	)
	if err != nil {
		_ = file.Close()

		return nil, fmt.Errorf("mmap file: %w", err)
	}

	if err := file.Close(); err != nil {
		_ = unix.Munmap(mapped)

		return nil, fmt.Errorf("close file: %w", err)
	}

	s := &sources{
		original: mapped,
		unmap: func() error {
			return unix.Munmap(mapped)
		},
	}

	pieces := makePieces(sourceOriginal, mapped)

	return &Buffer{
		sources: s,
		tree:    newPieceTree(s, pieces),
	}, nil
}
