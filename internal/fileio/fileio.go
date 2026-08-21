package fileio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"

	"github.com/i-walk-away/daun/internal/buffer"
)

const defaultFileMode os.FileMode = 0o644

// Open opens an existing regular file as a buffer.
//
// The original file contents are memory-mapped read-only instead of being
// copied into the Go heap. The returned buffer owns the mapping and must
// eventually be closed.
func Open(path string) (*buffer.Buffer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}

	info, err := file.Stat()
	if err != nil {
		_ = file.Close()

		return nil, fmt.Errorf("stat %q: %w", path, err)
	}

	if !info.Mode().IsRegular() {
		_ = file.Close()

		return nil, fmt.Errorf("path is not a regular file: %s", path)
	}

	if info.Size() == 0 {
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("close empty file: %w", err)
		}

		return buffer.NewBuffer(), nil
	}

	if info.Size() > int64(maxInt()) {
		_ = file.Close()

		return nil, fmt.Errorf(
			"file %q is too large: %d bytes",
			path,
			info.Size(),
		)
	}

	mapped, err := unix.Mmap(
		int(file.Fd()),
		0,
		int(info.Size()),
		unix.PROT_READ,
		unix.MAP_PRIVATE,
	)
	if err != nil {
		_ = file.Close()

		return nil, fmt.Errorf("mmap %q: %w", path, err)
	}

	if err := file.Close(); err != nil {
		_ = unix.Munmap(mapped)

		return nil, fmt.Errorf("close %q: %w", path, err)
	}

	return buffer.NewBufferFromReadOnlySource(
		mapped,
		func() error {
			return unix.Munmap(mapped)
		},
	), nil
}

// OpenOrCreate opens an existing regular file or creates an empty file.
//
// When the file does not exist, a new empty buffer is returned and the file
// is created immediately on disk.
func OpenOrCreate(path string) (*buffer.Buffer, error) {
	_, err := os.Stat(path)

	switch {
	case err == nil:
		return Open(path)

	case !os.IsNotExist(err):
		return nil, fmt.Errorf("stat %q: %w", path, err)
	}

	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_EXCL|os.O_WRONLY,
		defaultFileMode,
	)
	if err != nil {
		return nil, fmt.Errorf("create %q: %w", path, err)
	}

	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close %q: %w", path, err)
	}

	return buffer.NewBuffer(), nil
}

// Save atomically writes src to path.
//
// The document is first streamed into a temporary file in the same directory,
// synced to disk, and then renamed over the destination. The existing file
// permissions are preserved when possible.
func Save(path string, src io.WriterTo) error {
	if path == "" {
		return fmt.Errorf("file path is empty")
	}

	dir := filepath.Dir(path)
	base := filepath.Base(path)

	mode, err := fileMode(path)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(
		dir,
		"."+base+".daun-*",
	)
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}

	tmpPath := tmp.Name()

	cleanup := func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}

	defer cleanup()

	if err := tmp.Chmod(mode); err != nil {
		return fmt.Errorf("set temporary file permissions: %w", err)
	}

	if _, err := src.WriteTo(tmp); err != nil {
		return fmt.Errorf("write temporary file: %w", err)
	}

	if err := tmp.Sync(); err != nil {
		return fmt.Errorf("sync temporary file: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("replace %q: %w", path, err)
	}

	if err := syncDirectory(dir); err != nil {
		return fmt.Errorf("sync directory: %w", err)
	}

	return nil
}

func fileMode(path string) (os.FileMode, error) {
	info, err := os.Stat(path)

	switch {
	case err == nil:
		if !info.Mode().IsRegular() {
			return 0, fmt.Errorf("path is not a regular file: %s", path)
		}

		return info.Mode().Perm(), nil

	case os.IsNotExist(err):
		return defaultFileMode, nil

	default:
		return 0, fmt.Errorf("stat %q: %w", path, err)
	}
}

func syncDirectory(path string) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}

	defer dir.Close()

	return dir.Sync()
}

func maxInt() int {
	return int(^uint(0) >> 1)
}
