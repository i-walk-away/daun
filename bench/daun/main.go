package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/i-walk-away/daun/internal/buffer"
)

func main() {
	filePath := flag.String("file", "", "path to benchmark file")
	mode := flag.String("mode", "open", "benchmark mode: open, lookup, insert")
	iterations := flag.Int("iterations", 10_000, "number of operations")

	flag.Parse()

	if *filePath == "" {
		fatalf("file is required")
	}

	data, err := os.ReadFile(*filePath)
	if err != nil {
		fatalf("read file: %v", err)
	}

	fileSize := len(data)

	start := time.Now()

	buf := buffer.NewBufferFromText(string(data))

	openDuration := time.Since(start)

	fmt.Printf("file_bytes=%d\n", fileSize)
	fmt.Printf("line_count=%d\n", buf.LineCount())
	fmt.Printf("open_ns=%d\n", openDuration.Nanoseconds())

	switch *mode {
	case "open":
		return

	case "lookup":
		benchmarkLookup(buf, *iterations)

	case "insert":
		benchmarkInsert(buf, *iterations)

	default:
		fatalf("unknown mode %q", *mode)
	}
}

func benchmarkLookup(buf *buffer.Buffer, iterations int) {
	lineCount := buf.LineCount()

	if lineCount == 0 {
		return
	}

	start := time.Now()

	checksum := 0

	for i := 0; i < iterations; i++ {
		line := (i * 7919) % lineCount

		checksum += buf.LineLength(line)

		if i%100 == 0 {
			checksum += len(buf.Line(line))
		}
	}

	elapsed := time.Since(start)

	fmt.Printf("operation=lookup\n")
	fmt.Printf("iterations=%d\n", iterations)
	fmt.Printf("elapsed_ns=%d\n", elapsed.Nanoseconds())
	fmt.Printf(
		"ns_per_op=%.2f\n",
		float64(elapsed.Nanoseconds())/float64(iterations),
	)
	fmt.Printf("checksum=%d\n", checksum)
}

func benchmarkInsert(buf *buffer.Buffer, iterations int) {
	line := buf.LineCount() / 2
	column := buf.LineLength(line) / 2

	start := time.Now()

	for i := 0; i < iterations; i++ {
		buf.Insert(line, column, 'x')
		column++
	}

	elapsed := time.Since(start)

	fmt.Printf("operation=insert\n")
	fmt.Printf("iterations=%d\n", iterations)
	fmt.Printf("elapsed_ns=%d\n", elapsed.Nanoseconds())
	fmt.Printf(
		"ns_per_op=%.2f\n",
		float64(elapsed.Nanoseconds())/float64(iterations),
	)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
