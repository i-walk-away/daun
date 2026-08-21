package clipboard

// Clipboard provides access to the system clipboard.
type Clipboard interface {
	Copy(text string) error
	Paste() (string, error)
}
