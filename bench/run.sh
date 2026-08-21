#!/usr/bin/env bash

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

DATA_DIR="$ROOT/bench/data"
RESULTS_DIR="$ROOT/bench/results"
BIN_DIR="$ROOT/bench/.bin"

FILE_SIZE="${FILE_SIZE:-100M}"
ITERATIONS="${ITERATIONS:-10000}"
MEMORY_LIMIT="${MEMORY_LIMIT:-4G}"

FILE="$DATA_DIR/test-${FILE_SIZE}.txt"
RESULTS="$RESULTS_DIR/results.csv"

mkdir -p "$DATA_DIR"
mkdir -p "$RESULTS_DIR"
mkdir -p "$BIN_DIR"

require_command() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "error: required command '$1' not found" >&2
		exit 1
	fi
}

require_command go
require_command emacs
require_command systemd-run
require_command /usr/bin/time

generate_file() {
	local bytes

	case "$FILE_SIZE" in
		10M)
			bytes=10000000
			;;

		100M)
			bytes=100000000
			;;

		500M)
			bytes=500000000
			;;

		1G)
			bytes=1000000000
			;;

		*)
			echo "error: unsupported FILE_SIZE=$FILE_SIZE" >&2
			echo "supported values: 10M, 100M, 500M, 1G" >&2
			exit 1
			;;
	esac

	echo "Generating $FILE_SIZE benchmark file..."

	# 80 ASCII characters + newline per line.
	local line
	line="$(printf '%*s' 79 '' | tr ' ' 'a')"

	yes "$line" |
		head -c "$bytes" > "$FILE"

	echo
	echo "Generated:"
	wc -c -l "$FILE"
}

if [[ ! -f "$FILE" ]]; then
	generate_file
fi

echo "Building Daun benchmark adapter..."

go build \
	-o "$BIN_DIR/bench-daun" \
	./bench/daun

if [[ ! -f "$RESULTS" ]]; then
	echo "editor,size,mode,wall_seconds,max_rss_kb,cpu_percent" > "$RESULTS"
else
	: > "$RESULTS"
	echo "editor,size,mode,wall_seconds,max_rss_kb,cpu_percent" > "$RESULTS"
fi

run_benchmark() {
	local editor="$1"
	local mode="$2"
	shift 2

	local time_file
	local output_file

	time_file="$(mktemp)"
	output_file="$(mktemp)"

	echo
	echo "=== $editor / $mode / $FILE_SIZE ==="

	/usr/bin/time \
		-f '%e,%M,%P' \
		-o "$time_file" \
		systemd-run \
		--user \
		--scope \
		--quiet \
		--expand-environment=no \
		-p "MemoryMax=$MEMORY_LIMIT" \
		-p "MemorySwapMax=0" \
		"$@" >"$output_file"

	local wall
	local rss
	local cpu

	IFS=',' read -r wall rss cpu < "$time_file"

	echo "wall: ${wall}s"
	echo "rss:  ${rss} KB"
	echo "cpu:  ${cpu}"

	cat "$output_file"

	echo "$editor,$FILE_SIZE,$mode,$wall,$rss,$cpu" >> "$RESULTS"

	rm -f "$time_file"
	rm -f "$output_file"
}

for mode in open lookup insert; do
	run_benchmark \
		daun \
		"$mode" \
		"$BIN_DIR/bench-daun" \
		-file "$FILE" \
		-mode "$mode" \
		-iterations "$ITERATIONS"
done

for mode in open lookup insert; do
	run_benchmark \
		emacs \
		"$mode" \
		emacs \
		--batch \
		-Q \
		-l "$ROOT/bench/emacs.el" \
		-- \
		"$FILE" \
		"$mode" \
		"$ITERATIONS"
done

echo
echo "========================================"
echo "Benchmark finished"
echo "========================================"
echo
cat "$RESULTS"
echo
echo "Results saved to:"
echo "$RESULTS"