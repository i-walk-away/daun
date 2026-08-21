package buffer

type sourceKind uint8

const (
	sourceOriginal sourceKind = iota
	sourceAdd
)

type sources struct {
	original []byte
	add      []byte
	release  func() error
}

func (s *sources) bytes(source sourceKind, start, end int) []byte {
	switch source {
	case sourceOriginal:
		return s.original[start:end]

	case sourceAdd:
		return s.add[start:end]

	default:
		panic("buffer: unknown source")
	}
}

func (s *sources) close() error {
	if s.release == nil {
		return nil
	}

	return s.release()
}
