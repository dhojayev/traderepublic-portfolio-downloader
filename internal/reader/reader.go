//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=reader.go -destination reader_mock.go -package=reader

package reader

type Interface interface {
	Read(dataType string, req Request) (ResponseInterface, error)
}
