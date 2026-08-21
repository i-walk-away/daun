package clipboard

import (
	"fmt"
	"os/exec"
)

// XClip implements Clipboard using the xclip command.
type XClip struct{}

// NewXClip creates an XClip clipboard backend.
func NewXClip() *XClip {
	return &XClip{}
}

// Copy copies text to the X11 clipboard.
func (XClip) Copy(text string) error {
	cmd := exec.Command("xclip", "-selection", "clipboard")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("open xclip stdin: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start xclip: %w", err)
	}

	if _, err := stdin.Write([]byte(text)); err != nil {
		_ = stdin.Close()
		_ = cmd.Wait()

		return fmt.Errorf("write clipboard data: %w", err)
	}

	if err := stdin.Close(); err != nil {
		_ = cmd.Wait()

		return fmt.Errorf("close clipboard stdin: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("xclip failed: %w", err)
	}

	return nil
}

// Paste returns the current X11 clipboard contents.
func (XClip) Paste() (string, error) {
	out, err := exec.Command(
		"xclip",
		"-selection",
		"clipboard",
		"-o",
	).Output()
	if err != nil {
		return "", fmt.Errorf("read clipboard: %w", err)
	}

	return string(out), nil
}
