package filesystem

type OutputData struct {
	data []byte
}

func NewOutputData(data []byte) OutputData {
	return OutputData{data: data}
}

func (d OutputData) Data() []byte {
	return d.data
}
