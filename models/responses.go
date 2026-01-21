package models

// APIResponse is the standard JSON envelope for every HTTP response.
// Status mirrors the HTTP status text; Error is optional detail for failures.
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// BuildResponse constructs a consistent APIResponse body for handlers.
func BuildResponse(status string, message string, data interface{}, err string) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
		Error:   err,
	}
}
