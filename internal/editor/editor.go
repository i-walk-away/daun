package editor

type Editor struct {
	buffer Buffer
	cursor Position
}

type Position struct {
	Line int
	Char int
}

type Buffer struct {
	lines []string
}

func (e *Editor) Insert(r rune) {
}

// Backspace removes the rune right before the caret.
func (e *Editor) Backspace() {
}

// MoveLeft moves the caret left by 1 character.
func (e *Editor) MoveLeft() {
}

// MoveRight moves the caret right by 1 character.
func (e *Editor) MoveRight() {
}

// MoveUp moves the caret up by 1 line.
func (e *Editor) MoveUp() {
}

// MoveDown moves the caret down by 1 line.
func (e *Editor) MoveDown() {
}

// Enter inserts a newline at the current position.
func (e *Editor) Enter() {
}
