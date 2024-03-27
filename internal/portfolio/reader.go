package portfolio

type ReaderInterface interface {
	Read(dataType string, data map[string]any) (OutputDataInterface, error)
}
