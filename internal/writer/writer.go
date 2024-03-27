package writer

type Interface interface {
	Bytes(dir string, data []byte) error
}

type NilWriter struct{}

func NewNilWriter() NilWriter {
	return NilWriter{}
}

func (w NilWriter) Bytes(_ string, _ []byte) error {
	return nil
}
