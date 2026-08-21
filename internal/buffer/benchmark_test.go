package buffer

import (
	"math/rand"
	"strings"
	"testing"
)

func BenchmarkBufferInsertAtPosition(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"1MB", 1 << 20},
		{"10MB", 10 << 20},
		{"100MB", 100 << 20},
		{"1GB", 1 << 30},
	}

	for _, tt := range sizes {
		b.Run(tt.name, func(b *testing.B) {
			buf := NewBufferFromText(strings.Repeat("a", tt.size))

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				buf.Insert(0, 0, 'x')
			}
		})
	}
}

func BenchmarkBufferInsertAtEnd(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"1MB", 1 << 20},
		{"10MB", 10 << 20},
		{"100MB", 100 << 20},
	}

	for _, tt := range sizes {
		b.Run(tt.name, func(b *testing.B) {
			buf := NewBufferFromText(strings.Repeat("a", tt.size))
			column := buf.LineLength(0)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				buf.Insert(0, column, 'x')
			}
		})
	}
}

func BenchmarkBufferInsertAtMiddle(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"1MB", 1 << 20},
		{"10MB", 10 << 20},
		{"100MB", 100 << 20},
	}

	for _, tt := range sizes {
		b.Run(tt.name, func(b *testing.B) {
			buf := NewBufferFromText(strings.Repeat("a", tt.size))
			column := buf.LineLength(0) / 2

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				buf.Insert(0, column, 'x')
			}
		})
	}
}

func BenchmarkBufferInsertText(b *testing.B) {
	base := strings.Repeat("hello world\n", 20_000)

	sizes := []struct {
		name string
		size int
	}{
		{"1KB", 1 << 10},
		{"10KB", 10 << 10},
		{"100KB", 100 << 10},
	}

	for _, tt := range sizes {
		b.Run(tt.name, func(b *testing.B) {
			text := base[:min(len(base), tt.size)]

			if index := strings.LastIndexByte(text, '\n'); index >= 0 {
				text = text[:index+1]
			}

			b.ReportAllocs()
			b.SetBytes(int64(len(text)))

			for i := 0; i < b.N; i++ {
				buf := NewBuffer()
				buf.InsertText(0, 0, text)
			}
		})
	}
}

func BenchmarkBufferDeleteRange(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"1MB", 1 << 20},
		{"10MB", 10 << 20},
		{"100MB", 100 << 20},
	}

	for _, tt := range sizes {
		b.Run(tt.name, func(b *testing.B) {
			initial := strings.Repeat("a", tt.size)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				buf := NewBufferFromText(initial)

				start := tt.size / 3
				end := start + 1024

				buf.DeleteRange(0, start, end)
			}
		})
	}
}

func BenchmarkBufferTextRange(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"1MB", 1 << 20},
		{"10MB", 10 << 20},
		{"100MB", 100 << 20},
	}

	for _, tt := range sizes {
		b.Run(tt.name, func(b *testing.B) {
			buf := NewBufferFromText(strings.Repeat("a", tt.size))

			start := tt.size / 3
			end := start + 64*1024

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = buf.TextRange(0, start, 0, end)
			}
		})
	}
}

func BenchmarkBufferLineLookup(b *testing.B) {
	const lineCount = 1_000_000

	text := strings.Repeat("hello world\n", lineCount)
	buf := NewBufferFromText(text)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		line := i % lineCount
		_ = buf.Line(line)
	}
}

func BenchmarkBufferLineLength(b *testing.B) {
	const lineCount = 1_000_000

	text := strings.Repeat("hello world\n", lineCount)
	buf := NewBufferFromText(text)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		line := i % lineCount
		_ = buf.LineLength(line)
	}
}

func BenchmarkBufferRandomEdits(b *testing.B) {
	const initialSize = 100 * 1024 * 1024

	buf := NewBufferFromText(strings.Repeat("a", initialSize))
	rng := rand.New(rand.NewSource(42))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		size := buf.LineLength(0)

		if size < 2 {
			buf.Insert(0, 0, 'x')
			continue
		}

		column := rng.Intn(size)

		switch i % 3 {
		case 0:
			buf.Insert(0, column, 'x')

		case 1:
			buf.Delete(0, column)

		case 2:
			end := min(column+64, size)
			buf.DeleteRange(0, column, end)
		}
	}
}

func BenchmarkBufferRandomInsertPositions(b *testing.B) {
	const documentSize = 100 * 1024 * 1024

	buf := NewBufferFromText(strings.Repeat("a", documentSize))
	rng := rand.New(rand.NewSource(42))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		column := rng.Intn(buf.LineLength(0))
		buf.Insert(0, column, 'x')
	}
}

func BenchmarkBufferMultilineEdits(b *testing.B) {
	initial := strings.Repeat("hello world\n", 100_000)
	text := "inserted line\nwith multiple\nlines\n"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := NewBufferFromText(initial)

		buf.InsertText(50_000, 3, text)
		buf.DeleteRangeLines(
			25_000,
			2,
			25_010,
			5,
		)
	}
}
