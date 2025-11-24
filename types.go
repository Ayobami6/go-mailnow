package mailnow

// EmailRequest represents an email sending request
type EmailRequest struct {
	From        string       `json:"from"`
	To          string       `json:"to"`
	Subject     string       `json:"subject"`
	HTML        string       `json:"html"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Filename    string `json:"filename"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
}

// EmailResponse represents a successful email sending response
type EmailResponse struct {
	Data       Data   `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
}
type Data struct {
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
