package mailnow

// EmailRequest represents an email sending request
type EmailRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
}

// EmailResponse represents a successful email sending response
type EmailResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error struct {
		Code    string                 `json:"code"`
		Message string                 `json:"message"`
		Details map[string]interface{} `json:"details,omitempty"`
	} `json:"error"`
}
