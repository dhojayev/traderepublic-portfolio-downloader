package reader

// JSONResponse represents JSON response unmarshalled in bytes.
type JSONResponse struct {
	data []byte
}

func NewJSONResponse(data []byte) JSONResponse {
	return JSONResponse{data: data}
}

func (r JSONResponse) Data() []byte {
	return r.data
}
