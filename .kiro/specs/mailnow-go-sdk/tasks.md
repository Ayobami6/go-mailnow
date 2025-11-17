# Implementation Plan

- [x] 1. Set up project structure and Go module









  - Create directory structure for the Go SDK
  - Initialize go.mod with module path and Go 1.21 requirement
  - Create basic package structure (client.go, types.go, errors.go, validation.go, http.go, constants.go)
  - Add LICENSE file
  - _Requirements: 6.2, 6.3_

- [x] 2. Implement error types and error handling




  - Create base Error struct with Message and Err fields
  - Implement Error() and Unwrap() methods for base Error type
  - Create ValidationError, AuthError, RateLimitError, ServerError, and ConnectionError types
  - Add constructor functions for each error type
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_

- [x] 2.1 Write unit tests for error types











  - Test error creation and message formatting
  - Test error unwrapping with errors.Unwrap()
  - Test type assertions with errors.As()
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_

- [x] 3. Implement constants and configuration




  - Define API base URL, version, and endpoint constants
  - Define timeout constants
  - Define API key prefix constants
  - _Requirements: 2.2, 8.2_

- [x] 4. Implement validation functions





  - Create validateAPIKey() function to check for empty string and valid prefix format
  - Create validateEmailAddress() function using regex pattern
  - Create validateEmailRequest() function to validate all email parameters
  - _Requirements: 1.4, 5.1, 5.2, 5.3, 5.4, 5.5_

- [x] 4.1 Write unit tests for validation functions






  - Test validateAPIKey() with empty, invalid format, and valid API keys
  - Test validateEmailAddress() with invalid and valid email formats
  - Test validateEmailRequest() with missing, empty, and valid parameters
  - Use table-driven tests for comprehensive coverage
  - _Requirements: 1.4, 5.1, 5.2, 5.3, 5.4, 5.5_

- [x] 5. Implement request and response types




  - Create EmailRequest struct with json tags for from, to, subject, html fields
  - Create EmailResponse struct with json tags for success, message_id, status fields
  - Create ErrorResponse struct for API error responses
  - _Requirements: 2.1, 3.2_

- [x] 6. Implement HTTP client functions





  - Create makeRequest() function to build and send HTTP requests with proper headers
  - Implement JSON encoding for request body
  - Add X-API-Key and Content-Type headers
  - Create handleResponse() function to process HTTP responses
  - Implement status code to error type mapping (400→ValidationError, 401→AuthError, 429→RateLimitError, 5xx→ServerError)
  - Parse JSON error responses and extract error messages
  - _Requirements: 2.2, 2.3, 2.4, 2.5, 4.1, 4.2, 4.3, 4.4, 4.5_

- [x] 6.1 Write unit tests for HTTP functions






  - Test makeRequest() with proper header construction
  - Test handleResponse() for successful responses (2xx)
  - Test error mapping for each status code (400, 401, 429, 5xx)
  - Test connection error handling
  - Use httptest package to mock HTTP server responses
  - _Requirements: 2.2, 2.3, 2.4, 2.5, 4.1, 4.2, 4.3, 4.4, 4.5_

- [x] 7. Implement Client struct and constructor





  - Create Client struct with apiKey, httpClient, and baseURL fields
  - Implement NewClient() constructor function
  - Validate API key in NewClient() using validateAPIKey()
  - Initialize http.Client with timeout configuration
  - Return ValidationError if API key is invalid
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 8.1, 8.2, 8.3, 8.4_

- [x] 7.1 Write unit tests for Client constructor






  - Test NewClient() with empty API key (should return ValidationError)
  - Test NewClient() with invalid API key format (should return ValidationError)
  - Test NewClient() with valid mn_live_* API key (should succeed)
  - Test NewClient() with valid mn_test_* API key (should succeed)
  - _Requirements: 1.1, 1.2, 1.3, 1.4_

- [x] 8. Implement SendEmail method





  - Create SendEmail() method on Client struct
  - Accept context.Context and *EmailRequest parameters
  - Validate email request using validateEmailRequest()
  - Build full URL using baseURL and EmailSendEndpoint constant
  - Call makeRequest() to send HTTP POST request
  - Call handleResponse() to process response
  - Parse successful response JSON into EmailResponse struct
  - Return EmailResponse and nil error on success
  - Return appropriate error types on failure
  - _Requirements: 2.1, 2.2, 3.1, 3.2, 3.3, 5.1, 5.2, 5.3, 5.4, 5.5, 9.1, 9.2_

- [-] 8.1 Write unit tests for SendEmail method



  - Test SendEmail() with invalid request parameters (should return ValidationError)
  - Test SendEmail() with successful API response (should return EmailResponse)
  - Test SendEmail() with authentication error (should return AuthError)
  - Test SendEmail() with rate limit error (should return RateLimitError)
  - Test SendEmail() with server error (should return ServerError)
  - Test SendEmail() with context cancellation (should return error)
  - Use httptest to mock API responses
  - _Requirements: 2.1, 2.2, 3.1, 3.2, 3.3, 4.1, 4.2, 4.3, 4.4, 5.1, 5.2, 5.3, 5.4, 5.5, 9.2_

- [x] 9. Create comprehensive README documentation








  - Write installation instructions using go get
  - Add getting started section with basic usage example
  - Document API key format and where to obtain keys
  - Include complete code example with error handling
  - Document all error types and when they occur
  - Add context usage examples (timeout, cancellation)
  - Include link to Mailnow API documentation
  - Add badges for Go version, build status, and coverage
  - _Requirements: 6.4, 7.1, 7.2, 7.4_

- [x] 10. Add godoc comments to all exported types




  - Add package-level godoc comment explaining SDK purpose
  - Document Client struct with usage description
  - Document NewClient() function with parameters, returns, and errors
  - Document SendEmail() method with parameters, returns, and errors
  - Document EmailRequest and EmailResponse structs with field descriptions
  - Document all error types with when they are returned
  - Include example code in godoc comments where helpful
  - _Requirements: 7.3, 7.4, 9.4_

- [ ] 11. Set up GitHub Actions CI/CD workflow
  - Create .github/workflows/test.yml file
  - Configure workflow to run on push and pull request events
  - Set up matrix testing for Go versions (1.21, 1.22, latest)
  - Set up matrix testing for platforms (ubuntu, macos, windows)
  - Add step to run go test with coverage
  - Add step to run go vet for static analysis
  - Add step to check code formatting with gofmt
  - Add step to fail build if coverage is below 80%
  - _Requirements: 6.2, 9.5_

- [ ]* 12. Create integration tests
  - Create integration_test.go file with build tag for integration tests
  - Write test for successful email send with valid API key
  - Write test for authentication failure with invalid API key
  - Write test for validation errors with invalid email parameters
  - Write test for context timeout handling
  - Document how to run integration tests with mn_test_* API key
  - _Requirements: 2.1, 2.2, 3.1, 4.1, 5.5, 9.2_
