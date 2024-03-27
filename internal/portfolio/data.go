package portfolio

type OutputDataInterface interface {
	Data() []byte
}

type InputDataMap map[string]any
