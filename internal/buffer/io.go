package buffer

import (
	"fmt"
	"io"
)

// WriteTo writes the complete document to w.
//
// The document is streamed directly from the piece table without creating
// a second in-memory copy of the document.
//
// WriteTo implements io.WriterTo.
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	if b == nil || b.tree == nil {
		return 0, nil
	}

	var written int64

	err := b.tree.writeTo(func(data []byte) error {
		n, err := w.Write(data)

		written += int64(n)

		if err != nil {
			return err
		}

		if n != len(data) {
			return io.ErrShortWrite
		}

		return nil
	})

	if err != nil {
		return written, fmt.Errorf("write buffer: %w", err)
	}

	return written, nil
}

func (t *pieceTree) writeTo(write func([]byte) error) error {
	leaf := t.root

	for !leaf.leaf {
		leaf = leaf.children[0]
	}

	for ; leaf != nil; leaf = leaf.next {
		for _, p := range leaf.pieces {
			data := t.sources.bytes(p.source, p.start, p.end)

			if err := write(data); err != nil {
				return err
			}
		}
	}

	return nil
}
