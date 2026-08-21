# Daun

Daun is a simlple TUI text editor written in Go. I couldn't configure Emacs and decided to make my own thing.

### Buffer storage

The text itself is stored in the `Buffer`. The buffer uses a piece table backed by a B+ tree as its document storage
engine.

The piece table keeps the original document and newly inserted text in append-only storage. The B+ tree indexes pieces
and maintains aggregate metrics for each subtree, including byte count, rune count, and newline count.

This allows document operations to remain largely independent of the total document size.

1. Cursor and line lookup use the tree's aggregate metrics instead of scanning the entire document
2. Insertions and deletions modify only the affected tree path and pieces
3. Inserted text is stored once and referenced by pieces rather than repeatedly copied

The storage layer is intentionally hidden behind the `Buffer` API, so the editor itself does not depend on the underlying
data structure. These piece table shenanigans are merely an implementation detail of the `Buffer`.

### Editor model

The `Editor` keeps cursor position, selection state, and undo/redo history independently from the buffer.

Selections are represented by an anchor position and the current cursor position. This makes character, word, and
multiline selection behave consistently regardless of the direction in which the selection was created.

Editing operations are expressed as range replacements. This gives insertion, deletion, paste, selection replacement, and
undo/redo a common primitive.

### TUI

The terminal UI is responsible only for input handling and presentation. Everything under the `./internal/tui/` was AI
generated.

The viewport supports large documents by rendering only visible lines and automatically following the cursor. Vertical
and horizontal scrolling are independent from the editor's document coordinates.

Unicode is handled in terms of both **rune positions** and **terminal display width**, so cursor positioning and
rendering can account for wide characters.

### File handling

Daun can open existing files and create new files if the requested path does not exist:

```bash
daun filename
```

Opening a file does not require converting the entire document into one large Go string. File-backed documents are loaded
into the same `Buffer`
representation used by in-memory documents, so editing semantics remain identical regardless of how the document was
opened.

Saving (ctrl + S) is performed through the `fileio` layer.

### Search

Daun supports two search modes.

`Ctrl+F` opens text search. If some text is selected before invoking text search, the selection is used as the initial
query. Matches are highlighted in the document and the search UI tracks the current match and total number of matches.

Pressing `Ctrl+F` again switches to global file search. File search enumerates paths under the user's home directory,
ranks matching paths by filename relevance, and allows the selected result to be opened directly from the editor.

Text search operates over the piece table without materializing the entire document into a contiguous string. File search
uses `ripgrep` for filesystem enumeration and keeps only a bounded set of ranked results in memory.

### Benchmarks

Daun contains two different kinds of benchmarks.

Microbenchmarks live alongside the storage implementation and are used to measure individual buffer operations and detect
regressions in the piece table and B+ tree implementation.

Comparative benchmarks are kept separate from the storage microbenchmarks. They are intended to compare Daun against
other editors under identical workloads, input files, resource limits, and verification conditions.

The goal of the comparative benchmark is not to produce a single universal
"fastest editor" number. It measures concrete workloads such as opening, navigation, editing, searching, and saving while
recording wall-clock time, CPU time, and peak resident memory.

### Clipboard

Clipboard access is implemented through a small interface rather than being tied to the TUI or editor.

On Linux/X11, Daun currently uses `xclip` for system clipboard integration. This avoids relying on terminal-specific
OSC52 clipboard support and allows the clipboard backend to be replaced independently. Apparently, this approach is
infuriatingly TERRIBLE, but hey, it works on my end.

### Performance

The current storage benchmarks are intended to validate the scaling behavior of the document representation. They are not
a claim that Daun is universally faster or more memory-efficient than mature editors such as Emacs. Comparative editor
benchmarks are WIP and will be maintained separately.

The storage engine was stress-tested on very large documents.

A benchmark measuring insertion into documents ranging from **1 MB to 1 GB** produced approximately:

```text
1 MB     ~210 ns/op
10 MB    ~203 ns/op
100 MB   ~199 ns/op
1 GB     ~201 ns/op
```

Insertion cost therefore remained nearly constant across a 1000× increase in document size.

For larger text insertions, the benchmark reached approximately **1.7 GB/s** throughput.

The project was also stress-tested with documents containing **billions of lines and hundreds of billions of
characters**.

Initial buffer data structure was implemented using string slices and gap buffers. Stress testing showed noticable CPU
growth caused by maintaining millions of individual line buffers and large number of memory allocations associated with
large documents.

After reimplementing buffer storage with piece tables and b+ trees, insertion cost became independent of the document
size.

These benchmarks are primarily used as regression tests for the storage architecture rather than as absolute performance
guarantees across different hardware and workloads, of course.
