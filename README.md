# Daun

Daun is a simlple TUI text editor written in Go. I couldn't configure Emacs and decided to make my own thing.

# Features

Daun currently supports:

* Opening existing files from the command line: `daun file.txt`
* If the requested path does not exist, it creates a file
* Saving files with Ctrl+S
* Arrow keys move the cursor one character at a time. Holding Ctrl allows to navigate 1 *word* at a time. Works with
  selection (shift) too
* Copy, cut, and paste
* Undo and redo
* Search text inside current file with Ctrl+F
* Pressing Ctrl+F twice invokes a global home directory search by filename. Input a query and daun will find best
  matching filenames in your ~/ directory
* If any text is selected before pressing ctrl+f, the current selection will be automatically inputed as the initial
  search query
* Unicode-aware cursor positioning and rendering

# Building

Build the binary with:

```bash
make build
```

Then execute it:

```bash
./daun [optional file path]
```

You can also install it into your local binaries directory:

```bash
make install
``` 

The default installation path is `~/.local/bin/daun`. This will allow you to run daun from any directory.

### Buffer storage

The text itself is stored in the `Buffer`. The buffer uses a piece table backed by a B+ tree as its document storage
engine.

The piece table keeps the original document and newly inserted text in append-only storage. The B+ tree indexes pieces
and maintains aggregate metrics for each subtree.

Each tree node maintains aggregate metrics including:

* byte count
* rune count
* newline count

These metrics allow the buffer to locate lines and positions without scanning the entire document.

The design also avoids repeatedly copying inserted text, inserted data is stored once and referenced by pieces.

The storage layer is hidden behind the `Buffer` API, so the editor itself does not depend on the underlying data
structure. These piece table shenanigans are merely an implementation detail of the `Buffer`.

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

### Search

Daun supports two search modes.

`Ctrl+F` opens text search. If some text is selected before invoking text search, the selection is used as the initial
query. Matches are highlighted in the document and the search UI tracks the current match and total number of matches.

Pressing `Ctrl+F` again switches to global file search. File search enumerates paths under the user's home directory,
ranks matching paths by filename relevance, and allows the selected result to be opened directly from the editor.

Text search operates over the piece table without materializing the entire document into a contiguous string. File search
uses `ripgrep` for filesystem enumeration and keeps only a bounded set of ranked results in memory.

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
