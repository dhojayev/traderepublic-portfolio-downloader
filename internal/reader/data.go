package reader

type ResponseInterface interface {
	Data() []byte
}

type Request map[string]any
