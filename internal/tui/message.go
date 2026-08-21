package tui

import "time"

// MessageLevel identifies the severity of a UI message.
type MessageLevel uint8

const (
	MessageInfo MessageLevel = iota
	MessageSuccess
	MessageWarning
	MessageError
)

// Message represents a transient message shown in the editor message panel.
type Message struct {
	Level     MessageLevel
	Text      string
	CreatedAt time.Time
}
