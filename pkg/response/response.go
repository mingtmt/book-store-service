package response

// MessageResponse is a shared message response format.
type MessageResponse struct {
	Message string `json:"message"`
}

// ErrorResponse is a shared error response format.
type ErrorResponse struct {
	Error string `json:"error"`
}
