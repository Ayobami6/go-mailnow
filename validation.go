package mailnow

import (
	"regexp"
	"strings"
)

// emailRegex is a regex pattern for validating email addresses
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateAPIKey validates the API key format
func ValidateAPIKey(apiKey string) error {
	if apiKey == "" {
		return NewValidationError("API key cannot be empty", nil)
	}

	if !strings.HasPrefix(apiKey, APIKeyPrefixLive) && !strings.HasPrefix(apiKey, APIKeyPrefixTest) {
		return NewValidationError("API key must start with 'mn_live_' or 'mn_test_'", nil)
	}

	return nil
}

// ValidateEmailAddress validates an email address format
func ValidateEmailAddress(email string) error {
	if email == "" {
		return NewValidationError("email address cannot be empty", nil)
	}

	if !emailRegex.MatchString(email) {
		return NewValidationError("invalid email address format: "+email, nil)
	}

	return nil
}

// ValidateEmailRequest validates all email request parameters
func ValidateEmailRequest(req *EmailRequest) error {
	if req == nil {
		return NewValidationError("email request cannot be nil", nil)
	}

	// Validate from address
	if req.From == "" {
		return NewValidationError("from address is required", nil)
	}
	if err := ValidateEmailAddress(req.From); err != nil {
		return NewValidationError("invalid from address", err)
	}

	// Validate to address
	if req.To == "" {
		return NewValidationError("to address is required", nil)
	}
	if err := ValidateEmailAddress(req.To); err != nil {
		return NewValidationError("invalid to address", err)
	}

	// Validate subject
	if req.Subject == "" {
		return NewValidationError("subject is required", nil)
	}

	// Validate HTML body
	if req.HTML == "" {
		return NewValidationError("HTML body is required", nil)
	}

	return nil
}
