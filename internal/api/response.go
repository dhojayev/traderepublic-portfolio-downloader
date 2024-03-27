package api

type LoginRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	Pin         string `json:"pin"`
}

type LoginResponse struct {
	ProcessID string `json:"processId,omitempty"`
}
