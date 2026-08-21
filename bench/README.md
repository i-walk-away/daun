# Daun vs Emacs benchmarks

This directory contains reproducible large-file benchmarks for Daun and Emacs.

The benchmark is intentionally headless. It measures document storage and editing behavior.

## Workloads

### Open

Loads the same file into the editor's document representation.

Measured externally:

- wall-clock time;
- peak RSS;
- CPU usage.

### Lookup

Performs deterministic pseudo-random position lookups across the document.

Daun uses its buffer line index.

Emacs moves to deterministic character positions and queries the current line.

### Insert

Moves to the middle of the document and inserts one character repeatedly.

The same number of insertions is performed for both editors.

## Safety

Benchmarks run inside a user systemd scope with a memory limit.

Default:

- Memory: 4 GiB
- Swap: disabled

This prevents a benchmark from consuming all system memory.

## Running

Default benchmark:

    make bench-emacs

Environment variables:

    FILE_SIZE=10M
    FILE_SIZE=100M
    FILE_SIZE=500M
    FILE_SIZE=1G

Number of operations:

    ITERATIONS=10000

Memory limit:

    MEMORY_LIMIT=4G

Example:

    FILE_SIZE=500M ITERATIONS=5000 make bench-emacs

Results are written to:

    bench/results/results.csv

## Methodology

These benchmarks compare document-engine behavior rather than interactive GUI performance.

They should therefore be interpreted as storage/editing benchmarks, not as a complete comparison of editor usability or
terminal UI performance.

Always record:

- CPU model;
- system memory;
- operating system;
- Go version;
- Emacs version;
- Daun commit;
- benchmark file size.