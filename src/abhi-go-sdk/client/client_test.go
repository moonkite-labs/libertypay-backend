package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"abhi-go-sdk/errors"
	"abhi-go-sdk/models"
)

func TestNew(t *testing.T) {
	config := &Config{
		BaseURL:  "https://test.example.com",
		Username: "test",
		Password: "pass",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Timeout: 30 * time.Second,
	}

	client := New(config)

	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
	if client.config != config {
		t.Error("Expected config to be set")
	}
	if client.authManager == nil {
		t.Error("Expected authManager to be initialized")
	}
	if client.validator == nil {
		t.Error("Expected validator to be initialized")
	}
}

func TestNewWithNilConfig(t *testing.T) {
	client := New(nil)

	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
	if client.config == nil {
		t.Error("Expected default config to be set")
	}
}

func TestMakeRequestSuccess(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Expected Content-Type header to be application/json")
		}
		if r.Header.Get("Authorization") == "" {
			t.Error("Expected Authorization header to be set")
		}

		// Return success response
		response := models.APIResponse{
			StatusCode: 200,
			Message:    "Success",
			Data:       map[string]string{"test": "value"},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with test config
	config := &Config{
		BaseURL:  server.URL,
		Username: "test",
		Password: "pass",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Timeout: 30 * time.Second,
	}
	client := New(config)

	// Create a simple test auth manager for testing
	testAuthManager := &AuthManager{
		config: config,
		token: "test-token",
		expiresAt: time.Now().Add(time.Hour),
		httpClient: config.HTTPClient,
	}
	client.authManager = testAuthManager

	ctx := context.Background()
	var result map[string]string
	err := client.makeRequest(ctx, "GET", "/test", nil, &result)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result["test"] != "value" {
		t.Error("Expected result to contain test data")
	}
}

func TestMakeRequestValidationError(t *testing.T) {
	config := DefaultConfig()
	client := New(config)
	
	// Create a simple test auth manager for testing
	testAuthManager := &AuthManager{
		config: config,
		token: "test-token",
		expiresAt: time.Now().Add(time.Hour),
		httpClient: config.HTTPClient,
	}
	client.authManager = testAuthManager

	// Test with invalid struct that has validation tags
	type TestRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	ctx := context.Background()
	err := client.makeRequest(ctx, "POST", "/test", TestRequest{Email: "invalid"}, nil)

	if err == nil {
		t.Fatal("Expected validation error")
	}

	if _, ok := err.(*errors.ValidationError); !ok {
		t.Error("Expected ValidationError type")
	}
}

func TestMakeRequestAPIError(t *testing.T) {
	// Create a test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorResponse := models.ErrorResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Details:    "Invalid parameters",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
	}))
	defer server.Close()

	config := &Config{
		BaseURL:    server.URL,
		Username:   "test",
		Password:   "pass",
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		Timeout:    30 * time.Second,
	}
	client := New(config)
	// Create a simple test auth manager for testing
	testAuthManager := &AuthManager{
		config: config,
		token: "test-token",
		expiresAt: time.Now().Add(time.Hour),
		httpClient: config.HTTPClient,
	}
	client.authManager = testAuthManager

	ctx := context.Background()
	err := client.makeRequest(ctx, "GET", "/test", nil, nil)

	if err == nil {
		t.Fatal("Expected API error")
	}

	if apiErr, ok := err.(*errors.APIError); ok {
		if apiErr.StatusCode != 400 {
			t.Errorf("Expected status code 400, got %d", apiErr.StatusCode)
		}
		if apiErr.Message != "Bad Request" {
			t.Errorf("Expected message 'Bad Request', got %s", apiErr.Message)
		}
	} else {
		t.Error("Expected APIError type")
	}
}

func TestHTTPMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := models.APIResponse{
			StatusCode: 200,
			Message:    "Success",
			Data:       map[string]string{"method": r.Method},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	config := &Config{
		BaseURL:    server.URL,
		Username:   "test",
		Password:   "pass",
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		Timeout:    30 * time.Second,
	}
	client := New(config)
	// Create a simple test auth manager for testing
	testAuthManager := &AuthManager{
		config: config,
		token: "test-token",
		expiresAt: time.Now().Add(time.Hour),
		httpClient: config.HTTPClient,
	}
	client.authManager = testAuthManager

	ctx := context.Background()

	// Test GET
	var getResult map[string]string
	err := client.GET(ctx, "/test", &getResult)
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if getResult["method"] != "GET" {
		t.Error("Expected GET method")
	}

	// Test POST with proper struct that has validation tags
	type TestBody struct {
		Data string `json:"data"`
	}
	var postResult map[string]string
	err = client.POST(ctx, "/test", TestBody{Data: "test"}, &postResult)
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	if postResult["method"] != "POST" {
		t.Error("Expected POST method")
	}

	// Test PUT
	var putResult map[string]string
	err = client.PUT(ctx, "/test", TestBody{Data: "test"}, &putResult)
	if err != nil {
		t.Fatalf("PUT failed: %v", err)
	}
	if putResult["method"] != "PUT" {
		t.Error("Expected PUT method")
	}

	// Test DELETE
	var deleteResult map[string]string
	err = client.DELETE(ctx, "/test", &deleteResult)
	if err != nil {
		t.Fatalf("DELETE failed: %v", err)
	}
	if deleteResult["method"] != "DELETE" {
		t.Error("Expected DELETE method")
	}
}

func TestSetRetryPolicy(t *testing.T) {
	config := DefaultConfig()
	client := New(config)

	originalTransport := client.httpClient.Transport

	client.SetRetryPolicy(3, 5*time.Second)

	if client.httpClient.Transport == originalTransport {
		t.Error("Expected transport to be wrapped with retry logic")
	}

	// Verify it's a retry transport
	if _, ok := client.httpClient.Transport.(*retryTransport); !ok {
		t.Error("Expected transport to be retryTransport")
	}
}

func TestRetryTransportSuccess(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	transport := &retryTransport{
		transport:  http.DefaultTransport,
		maxRetries: 3,
		retryDelay: 10 * time.Millisecond,
	}

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := transport.RoundTrip(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestRetryTransportWithRetries(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		w.Write([]byte("response"))
	}))
	defer server.Close()

	transport := &retryTransport{
		transport:  http.DefaultTransport,
		maxRetries: 3,
		retryDelay: 10 * time.Millisecond,
	}

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, err := transport.RoundTrip(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryTransportWithBodyRetries(t *testing.T) {
	attempts := 0
	var receivedBodies []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		
		body := make([]byte, 1024)
		n, _ := r.Body.Read(body)
		receivedBodies = append(receivedBodies, string(body[:n]))

		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	transport := &retryTransport{
		transport:  http.DefaultTransport,
		maxRetries: 3,
		retryDelay: 10 * time.Millisecond,
	}

	body := "test request body"
	req, _ := http.NewRequest("POST", server.URL, bytes.NewBufferString(body))
	resp, err := transport.RoundTrip(req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}

	// Verify body was sent correctly in all attempts
	for i, receivedBody := range receivedBodies {
		if receivedBody != body {
			t.Errorf("Attempt %d: expected body %q, got %q", i+1, body, receivedBody)
		}
	}
}

