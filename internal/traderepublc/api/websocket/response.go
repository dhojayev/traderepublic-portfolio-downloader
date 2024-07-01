package websocket

type ResponseErrors struct {
	Errors []ResponseError `json:"errors"`
}

type ResponseError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
