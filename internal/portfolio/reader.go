//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=reader.go -destination reader_mock.go -package=portfolio

package portfolio

type ReaderInterface interface {
	Read(dataType string, data map[string]any) (OutputDataInterface, error)
}
