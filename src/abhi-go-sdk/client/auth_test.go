package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"abhi-go-sdk/models"
	"github.com/golang-jwt/jwt/v4"
)

func TestNewAuthManager(t *testing.T) {
	config := &Config{
		Username: "test",
		Password: "pass",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	authManager := NewAuthManager(config)

	if authManager == nil {
		t.Fatal("Expected authManager to be non-nil")
	}
	if authManager.config != config {
		t.Error("Expected config to be set")
	}
	if authManager.httpClient != config.HTTPClient {
		t.Error("Expected httpClient to be set from config")
	}
}

func TestGetTokenFirstTime(t *testing.T) {
	// Create a test server for login
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/login" {
			t.Errorf("Expected path /auth/login, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify request body
		var loginReq models.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&loginReq)
		if err != nil {
			t.Fatalf("Failed to decode login request: %v", err)
		}
		if loginReq.Username != "test" || loginReq.Password != "pass" {
			t.Error("Expected correct username and password")
		}

		// Create a test JWT token
		token := createTestJWT(time.Now().Add(time.Hour))
		
		response := models.APIResponse{
			StatusCode: 200,
			Message:    "Success",
			Data: map[string]interface{}{
				"token": token,
			},
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
	}

	authManager := NewAuthManager(config)
	ctx := context.Background()

	token, err := authManager.GetToken(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if token == "" {
		t.Error("Expected non-empty token")
	}
}

func TestGetTokenCached(t *testing.T) {
	serverCallCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCallCount++
		
		token := createTestJWT(time.Now().Add(time.Hour))
		response := models.APIResponse{
			StatusCode: 200,
			Message:    "Success",
			Data: map[string]interface{}{
				"token": token,
			},
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
	}

	authManager := NewAuthManager(config)
	ctx := context.Background()

	// First call should hit the server
	token1, err := authManager.GetToken(ctx)
	if err != nil {
		t.Fatalf("Expected no error on first call, got %v", err)
	}

	// Second call should use cached token
	token2, err := authManager.GetToken(ctx)
	if err != nil {
		t.Fatalf("Expected no error on second call, got %v", err)
	}

	if token1 != token2 {
		t.Error("Expected same token from cache")
	}
	if serverCallCount != 1 {
		t.Errorf("Expected server to be called once, got %d calls", serverCallCount)
	}
}

func TestGetTokenExpired(t *testing.T) {
	serverCallCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCallCount++
		
		// First call returns expired token, second call returns valid token
		var expiry time.Time
		if serverCallCount == 1 {
			expiry = time.Now().Add(-time.Hour) // Expired
		} else {
			expiry = time.Now().Add(time.Hour) // Valid
		}
		
		token := createTestJWT(expiry)
		response := models.APIResponse{
			StatusCode: 200,
			Message:    "Success",
			Data: map[string]interface{}{
				"token": token,
			},
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
	}

	authManager := NewAuthManager(config)
	ctx := context.Background()

	// First call gets expired token
	_, err := authManager.GetToken(ctx)
	if err != nil {
		t.Fatalf("Expected no error on first call, got %v", err)
	}

	// Second call should refresh the token
	_, err = authManager.GetToken(ctx)
	if err != nil {
		t.Fatalf("Expected no error on second call, got %v", err)
	}

	if serverCallCount != 2 {
		t.Errorf("Expected server to be called twice for refresh, got %d calls", serverCallCount)
	}
}

func TestIsTokenValid(t *testing.T) {
	config := &Config{
		Username:   "test",
		Password:   "pass",
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
	authManager := NewAuthManager(config)

	// Test empty token
	if authManager.isTokenValid() {
		t.Error("Expected empty token to be invalid")
	}

	// Test expired token
	authManager.token = "some-token"
	authManager.expiresAt = time.Now().Add(-time.Hour)
	if authManager.isTokenValid() {
		t.Error("Expected expired token to be invalid")
	}

	// Test token expiring soon (within 5 minutes)
	authManager.expiresAt = time.Now().Add(3 * time.Minute)
	if authManager.isTokenValid() {
		t.Error("Expected token expiring soon to be invalid")
	}

	// Test valid token
	authManager.expiresAt = time.Now().Add(time.Hour)
	if !authManager.isTokenValid() {
		t.Error("Expected valid token to be valid")
	}
}

func TestParseTokenExpiration(t *testing.T) {
	authManager := &AuthManager{}

	// Test valid JWT token
	expiry := time.Now().Add(time.Hour)
	token := createTestJWT(expiry)

	parsedExpiry, err := authManager.parseTokenExpiration(token)
	if err != nil {
		t.Fatalf("Expected no error parsing valid token, got %v", err)
	}

	// Allow 1 second difference for test execution time
	if parsedExpiry.Unix() != expiry.Unix() {
		t.Errorf("Expected expiry %v, got %v", expiry.Unix(), parsedExpiry.Unix())
	}

	// Test invalid token
	_, err = authManager.parseTokenExpiration("invalid-token")
	if err == nil {
		t.Error("Expected error parsing invalid token")
	}
}

func TestClearToken(t *testing.T) {
	authManager := &AuthManager{
		token:     "some-token",
		expiresAt: time.Now().Add(time.Hour),
	}

	authManager.ClearToken()

	if authManager.token != "" {
		t.Error("Expected token to be cleared")
	}
	if !authManager.expiresAt.IsZero() {
		t.Error("Expected expiresAt to be cleared")
	}
}

func TestRefreshTokenServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	config := &Config{
		BaseURL:    server.URL,
		Username:   "test",
		Password:   "pass",
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}

	authManager := NewAuthManager(config)
	ctx := context.Background()

	_, err := authManager.GetToken(ctx)

	if err == nil {
		t.Error("Expected error for server error response")
	}
}

func TestRefreshTokenInvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return response without token
		response := models.APIResponse{
			StatusCode: 200,
			Message:    "Success",
			Data:       map[string]interface{}{},
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
	}

	authManager := NewAuthManager(config)
	ctx := context.Background()

	_, err := authManager.GetToken(ctx)

	if err == nil {
		t.Error("Expected error for response without token")
	}
}

// Helper function to create test JWT tokens
func createTestJWT(expiry time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expiry.Unix(),
		"iat": time.Now().Unix(),
		"sub": "test-user",
	})

	tokenString, _ := token.SignedString([]byte("test-secret"))
	return tokenString
}

// Benchmark tests
func BenchmarkGetTokenCached(b *testing.B) {
	authManager := &AuthManager{
		token:     "cached-token",
		expiresAt: time.Now().Add(time.Hour),
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := authManager.GetToken(ctx)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkIsTokenValid(b *testing.B) {
	authManager := &AuthManager{
		token:     "test-token",
		expiresAt: time.Now().Add(time.Hour),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = authManager.isTokenValid()
	}
}