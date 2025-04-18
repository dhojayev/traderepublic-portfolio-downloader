package api

// WebSocket response types

type WSListResponse struct {
	Items   []any                 `json:"items"`
	Cursors WSListCursorsResponse `json:"cursors"`
}

type WSListCursorsResponse struct {
	Before string `json:"before"`
	After  string `json:"after,omitempty"`
}
