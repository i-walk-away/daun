package buffer

import (
	"strings"
	"testing"
)

var (
	benchmarkSearchResult SearchResult
	benchmarkSearchLine   int
	benchmarkSearchColumn int
	benchmarkSearchOK     bool
)

func BenchmarkBufferFind(b *testing.B) {
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
			text := strings.Repeat("a", tt.size)

			matchOffset := tt.size / 2

			text = text[:matchOffset] +
				"TARGET" +
				text[matchOffset:]

			buf := NewBufferFromText(text)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				benchmarkSearchLine,
					benchmarkSearchColumn,
					benchmarkSearchOK = buf.Find(
					0,
					0,
					"TARGET",
				)
			}
		})
	}
}

func BenchmarkBufferFindPrevious(b *testing.B) {
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
			text := strings.Repeat("a", tt.size)

			matchOffset := tt.size / 2

			text = text[:matchOffset] +
				"TARGET" +
				text[matchOffset:]

			buf := NewBufferFromText(text)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				benchmarkSearchLine,
					benchmarkSearchColumn,
					benchmarkSearchOK = buf.FindPrevious(
					0,
					tt.size,
					"TARGET",
				)
			}
		})
	}
}

func BenchmarkBufferSearch(b *testing.B) {
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
			text := strings.Repeat("a", tt.size)

			matchOffset := tt.size / 2

			text = text[:matchOffset] +
				"TARGET" +
				text[matchOffset:]

			buf := NewBufferFromText(text)

			b.ReportAllocs()
			b.SetBytes(int64(tt.size))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				benchmarkSearchResult = buf.Search(
					0,
					0,
					"TARGET",
				)
			}
		})
	}
}

func BenchmarkBufferSearchManyMatches(b *testing.B) {
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
			const chunk = "TARGET\n"

			repetitions := tt.size / len(chunk)

			text := strings.Repeat(chunk, repetitions)

			buf := NewBufferFromText(text)

			b.ReportAllocs()
			b.SetBytes(int64(len(text)))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				benchmarkSearchResult = buf.Search(
					0,
					0,
					"TARGET",
				)
			}
		})
	}
}

func BenchmarkBufferFindAcrossPieces(b *testing.B) {
	const (
		prefix = "TARGET"
		size   = 1 << 20
	)

	text := strings.Repeat("a", size)

	buf := NewBufferFromText(text)

	// Force the searched text through edits so it is represented by multiple
	// pieces rather than one contiguous original-source slice.
	buf.InsertText(0, size/2, prefix)

	b.ReportAllocs()
	b.SetBytes(size)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkSearchLine,
			benchmarkSearchColumn,
			benchmarkSearchOK = buf.Find(
			0,
			size/2,
			prefix,
		)
	}
}
