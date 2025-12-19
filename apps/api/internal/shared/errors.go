package shared

type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

func NewError(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	}
}
