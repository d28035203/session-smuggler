package models

// APIResponse is the standard envelope for all HTTP responses.
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// BuildResponse constructs a consistent API response body.
func BuildResponse(status string, message string, data interface{}, err string) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
		Error:   err,
	}
}
