package writer

type Writer interface {
	Bytes(filename string, data []byte) error
}