package filesystem

type OutputData struct {
	data []byte
}

func (d OutputData) Data() []byte {
	return d.data
}
