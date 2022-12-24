package scopedata

type Reader interface {
	Read() ([]Value, error)
}

func NewReader() Reader {
	return &reader{}
}

type reader struct {
}

func (r *reader) Read() ([]Value, error) {
	return nil, nil
}
