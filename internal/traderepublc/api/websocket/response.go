package websocket

const (
	errorMsgAuthError         = "Authentication Error"
	errorMsgUnauthorizedError = "Unauthorized"
)

type ResponseErrors struct {
	Errors []ResponseError `json:"errors"`
}

type ResponseError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (e ResponseError) IsAuthError() bool {
	return e.ErrorMessage == errorMsgAuthError
}

func (e ResponseError) IsUnauthorizedError() bool {
	return e.ErrorMessage == errorMsgUnauthorizedError
}
