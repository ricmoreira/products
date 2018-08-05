package mresponse

type (
	// ErrorResponse represents the error model for App error response messages
	ErrorResponse struct {
		HttpCode int           `json:"-"`
		Code     string        `json:"code"`
		Response string        `json:"response"`
		Errors   []ErrorDetail `json:"errors,omitempty"`
	}

	// ErrorDetail represents a detailed error message applyed to the specific field or property that caused an error to the response
	ErrorDetail struct {
		Property string `json:"property"`
		Message  string `json:"message"`
	}
)
