package websocket

type CollectionResponse[T any] struct {
	Items   []T                       `json:"items"`
	Cursors CollectionCursorsResponse `json:"cursors"`
}

type CollectionCursorsResponse struct {
	Before string `json:"before"`
	After  string `json:"after,omitempty"`
}

type ResponseErrors struct {
	Errors []ResponseError `json:"errors"`
}

type ResponseError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
