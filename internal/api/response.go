package api

type LoginRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Pin         string `json:"pin"`
}

type LoginResponse struct {
	ProcessID string `json:"processId,omitempty"`
}

type WSListResponse struct {
	Items   []any                 `json:"items"`
	Cursors WSListCursorsResponse `json:"cursors"`
}

type WSListCursorsResponse struct {
	Before string `json:"before"`
	After  string `json:"after,omitempty"`
}
