package errors

import (
	"fmt"
	"net/http"
)

// APIError represents an error from the Abhi API
type APIError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	Endpoint   string `json:"endpoint,omitempty"`
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API Error [%d]: %s - %s", e.StatusCode, e.Message, e.Details)
	}
	return fmt.Sprintf("API Error [%d]: %s", e.StatusCode, e.Message)
}

// IsClientError returns true if the error is a 4xx client error
func (e *APIError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsServerError returns true if the error is a 5xx server error
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500
}

// IsUnauthorized returns true if the error is a 401 Unauthorized
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsForbidden returns true if the error is a 403 Forbidden
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == http.StatusForbidden
}

// IsNotFound returns true if the error is a 404 Not Found
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsConflict returns true if the error is a 409 Conflict
func (e *APIError) IsConflict() bool {
	return e.StatusCode == http.StatusConflict
}

// IsRateLimited returns true if the error is a 429 Too Many Requests
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// NewAPIError creates a new API error
func NewAPIError(statusCode int, message, details, endpoint string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
		Endpoint:   endpoint,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error for field '%s': %s", e.Field, e.Message)
}

// NetworkError represents a network-related error
type NetworkError struct {
	Operation string
	Err       error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("Network error during %s: %v", e.Operation, e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// AuthenticationError represents an authentication-related error
type AuthenticationError struct {
	Message string
	Err     error
}

func (e *AuthenticationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Authentication error: %s - %v", e.Message, e.Err)
	}
	return fmt.Sprintf("Authentication error: %s", e.Message)
}

func (e *AuthenticationError) Unwrap() error {
	return e.Err
}