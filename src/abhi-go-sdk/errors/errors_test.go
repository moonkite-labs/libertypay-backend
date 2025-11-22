package errors

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAPIError(t *testing.T) {
	err := &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Bad Request",
		Details:    "Invalid parameter",
		Endpoint:   "/test",
	}

	expected := "API Error [400]: Bad Request - Invalid parameter"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestAPIErrorWithoutDetails(t *testing.T) {
	err := &APIError{
		StatusCode: http.StatusNotFound,
		Message:    "Not Found",
		Endpoint:   "/test",
	}

	expected := "API Error [404]: Not Found"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestAPIErrorIsClientError(t *testing.T) {
	tests := []struct {
		statusCode    int
		expectedResult bool
	}{
		{http.StatusBadRequest, true},      // 400
		{http.StatusUnauthorized, true},    // 401
		{http.StatusForbidden, true},       // 403
		{http.StatusNotFound, true},        // 404
		{http.StatusConflict, true},        // 409
		{http.StatusTooManyRequests, true}, // 429
		{499, true},                        // Any 4xx
		{http.StatusOK, false},             // 200
		{http.StatusInternalServerError, false}, // 500
		{http.StatusBadGateway, false},     // 502
	}

	for _, test := range tests {
		err := &APIError{StatusCode: test.statusCode}
		result := err.IsClientError()
		if result != test.expectedResult {
			t.Errorf("StatusCode %d: expected IsClientError() to be %v, got %v",
				test.statusCode, test.expectedResult, result)
		}
	}
}

func TestAPIErrorIsServerError(t *testing.T) {
	tests := []struct {
		statusCode    int
		expectedResult bool
	}{
		{http.StatusInternalServerError, true}, // 500
		{http.StatusBadGateway, true},          // 502
		{http.StatusServiceUnavailable, true},  // 503
		{599, true},                            // Any 5xx
		{http.StatusOK, false},                 // 200
		{http.StatusBadRequest, false},         // 400
		{http.StatusNotFound, false},           // 404
	}

	for _, test := range tests {
		err := &APIError{StatusCode: test.statusCode}
		result := err.IsServerError()
		if result != test.expectedResult {
			t.Errorf("StatusCode %d: expected IsServerError() to be %v, got %v",
				test.statusCode, test.expectedResult, result)
		}
	}
}

func TestAPIErrorIsUnauthorized(t *testing.T) {
	tests := []struct {
		statusCode    int
		expectedResult bool
	}{
		{http.StatusUnauthorized, true}, // 401
		{http.StatusOK, false},          // 200
		{http.StatusForbidden, false},   // 403
		{http.StatusNotFound, false},    // 404
	}

	for _, test := range tests {
		err := &APIError{StatusCode: test.statusCode}
		result := err.IsUnauthorized()
		if result != test.expectedResult {
			t.Errorf("StatusCode %d: expected IsUnauthorized() to be %v, got %v",
				test.statusCode, test.expectedResult, result)
		}
	}
}

func TestAPIErrorIsForbidden(t *testing.T) {
	tests := []struct {
		statusCode    int
		expectedResult bool
	}{
		{http.StatusForbidden, true},    // 403
		{http.StatusOK, false},          // 200
		{http.StatusUnauthorized, false}, // 401
		{http.StatusNotFound, false},    // 404
	}

	for _, test := range tests {
		err := &APIError{StatusCode: test.statusCode}
		result := err.IsForbidden()
		if result != test.expectedResult {
			t.Errorf("StatusCode %d: expected IsForbidden() to be %v, got %v",
				test.statusCode, test.expectedResult, result)
		}
	}
}

func TestAPIErrorIsNotFound(t *testing.T) {
	tests := []struct {
		statusCode    int
		expectedResult bool
	}{
		{http.StatusNotFound, true},     // 404
		{http.StatusOK, false},          // 200
		{http.StatusUnauthorized, false}, // 401
		{http.StatusForbidden, false},   // 403
	}

	for _, test := range tests {
		err := &APIError{StatusCode: test.statusCode}
		result := err.IsNotFound()
		if result != test.expectedResult {
			t.Errorf("StatusCode %d: expected IsNotFound() to be %v, got %v",
				test.statusCode, test.expectedResult, result)
		}
	}
}

func TestAPIErrorIsConflict(t *testing.T) {
	tests := []struct {
		statusCode    int
		expectedResult bool
	}{
		{http.StatusConflict, true},     // 409
		{http.StatusOK, false},          // 200
		{http.StatusBadRequest, false},  // 400
		{http.StatusNotFound, false},    // 404
	}

	for _, test := range tests {
		err := &APIError{StatusCode: test.statusCode}
		result := err.IsConflict()
		if result != test.expectedResult {
			t.Errorf("StatusCode %d: expected IsConflict() to be %v, got %v",
				test.statusCode, test.expectedResult, result)
		}
	}
}

func TestAPIErrorIsRateLimited(t *testing.T) {
	tests := []struct {
		statusCode    int
		expectedResult bool
	}{
		{http.StatusTooManyRequests, true}, // 429
		{http.StatusOK, false},             // 200
		{http.StatusBadRequest, false},     // 400
		{http.StatusUnauthorized, false},   // 401
	}

	for _, test := range tests {
		err := &APIError{StatusCode: test.statusCode}
		result := err.IsRateLimited()
		if result != test.expectedResult {
			t.Errorf("StatusCode %d: expected IsRateLimited() to be %v, got %v",
				test.statusCode, test.expectedResult, result)
		}
	}
}

func TestNewAPIError(t *testing.T) {
	statusCode := http.StatusBadRequest
	message := "Bad Request"
	details := "Invalid parameter"
	endpoint := "/test"

	err := NewAPIError(statusCode, message, details, endpoint)

	if err.StatusCode != statusCode {
		t.Errorf("Expected StatusCode %d, got %d", statusCode, err.StatusCode)
	}
	if err.Message != message {
		t.Errorf("Expected Message '%s', got '%s'", message, err.Message)
	}
	if err.Details != details {
		t.Errorf("Expected Details '%s', got '%s'", details, err.Details)
	}
	if err.Endpoint != endpoint {
		t.Errorf("Expected Endpoint '%s', got '%s'", endpoint, err.Endpoint)
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "email",
		Message: "Invalid email format",
		Value:   "invalid-email",
	}

	expected := "Validation error for field 'email': Invalid email format"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestNetworkError(t *testing.T) {
	innerErr := fmt.Errorf("connection timeout")
	err := &NetworkError{
		Operation: "POST /api/test",
		Err:       innerErr,
	}

	expected := "Network error during POST /api/test: connection timeout"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}

	// Test Unwrap
	if err.Unwrap() != innerErr {
		t.Error("Expected Unwrap to return the inner error")
	}
}

func TestAuthenticationErrorWithInnerError(t *testing.T) {
	innerErr := fmt.Errorf("token expired")
	err := &AuthenticationError{
		Message: "Authentication failed",
		Err:     innerErr,
	}

	expected := "Authentication error: Authentication failed - token expired"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}

	// Test Unwrap
	if err.Unwrap() != innerErr {
		t.Error("Expected Unwrap to return the inner error")
	}
}

func TestAuthenticationErrorWithoutInnerError(t *testing.T) {
	err := &AuthenticationError{
		Message: "Invalid credentials",
	}

	expected := "Authentication error: Invalid credentials"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}

	// Test Unwrap with nil inner error
	if err.Unwrap() != nil {
		t.Error("Expected Unwrap to return nil when no inner error")
	}
}

// Benchmark tests
func BenchmarkAPIErrorError(b *testing.B) {
	err := &APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Bad Request",
		Details:    "Invalid parameter",
		Endpoint:   "/test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkAPIErrorIsClientError(b *testing.B) {
	err := &APIError{StatusCode: http.StatusBadRequest}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.IsClientError()
	}
}

func BenchmarkNewAPIError(b *testing.B) {
	statusCode := http.StatusBadRequest
	message := "Bad Request"
	details := "Invalid parameter"
	endpoint := "/test"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewAPIError(statusCode, message, details, endpoint)
	}
}