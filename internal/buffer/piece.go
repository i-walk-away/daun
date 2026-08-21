package buffer

import (
	"bytes"
	"unicode/utf8"
)

const maxPieceBytes = 32 * 1024

type metric struct {
	bytes    int
	runes    int
	newlines int
}

type piece struct {
	source  sourceKind
	start   int
	end     int
	metrics metric
}

func newPiece(source sourceKind, start, end int, data []byte) piece {
	return piece{
		source:  source,
		start:   start,
		end:     end,
		metrics: measure(data[start:end]),
	}
}

func measure(data []byte) metric {
	return metric{
		bytes:    len(data),
		runes:    utf8.RuneCount(data),
		newlines: bytes.Count(data, []byte{'\n'}),
	}
}

func addMetric(a, b metric) metric {
	return metric{
		bytes:    a.bytes + b.bytes,
		runes:    a.runes + b.runes,
		newlines: a.newlines + b.newlines,
	}
}

func (p piece) byteLen() int {
	return p.end - p.start
}

func canMerge(a, b piece) bool {
	return a.source == b.source &&
		a.end == b.start &&
		a.byteLen()+b.byteLen() <= maxPieceBytes
}

func mergePieces(a, b piece) piece {
	return piece{
		source:  a.source,
		start:   a.start,
		end:     b.end,
		metrics: addMetric(a.metrics, b.metrics),
	}
}
