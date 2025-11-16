package mailnow

import "time"

const (
	// APIBaseURL is the base URL for the Mailnow API
	APIBaseURL = "https://api.mailnow.xyz"

	// APIVersion is the API version
	APIVersion = "v1"

	// EmailSendEndpoint is the endpoint for sending emails
	EmailSendEndpoint = "/v1/email/send"

	// RequestTimeout is the default timeout for API requests
	RequestTimeout = 30 * time.Second

	// APIKeyPrefixLive is the prefix for live API keys
	APIKeyPrefixLive = "mn_live_"

	// APIKeyPrefixTest is the prefix for test API keys
	APIKeyPrefixTest = "mn_test_"
)
