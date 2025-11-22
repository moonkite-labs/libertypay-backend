package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"abhi-go-sdk/errors"
	"abhi-go-sdk/models"
	"github.com/go-playground/validator/v10"
	pkgerrors "github.com/pkg/errors"
)

// Client represents the Abhi API client
type Client struct {
	config            *Config
	authManager       *AuthManager
	httpClient        *http.Client
	validator         *validator.Validate
	rateLimiter       *RateLimiter
	credentialManager *CredentialManager
	requestSigner     *RequestSigner
}

// New creates a new Abhi API client
func New(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	client := &Client{
		config:      config,
		authManager: NewAuthManager(config),
		httpClient:  config.HTTPClient,
		validator:   validator.New(),
		rateLimiter: NewRateLimiter(config.RateLimit),
	}

	// Initialize security features
	if config.Security != nil {
		// Initialize credential manager if encryption is enabled
		if config.Security.EncryptCredentials && config.Security.EncryptionPassword != "" {
			client.credentialManager = NewCredentialManager(
				config.Security.EncryptionPassword,
				config.Security.CredentialStore,
			)
		}

		// Initialize request signer if enabled
		if config.Security.EnableRequestSigning && config.Security.SigningSecret != "" {
			client.requestSigner = NewRequestSigner(config.Security.SigningSecret)
		}
	}

	// Wrap HTTP client with middleware (rate limiting, signing)
	if client.httpClient != nil {
		transport := client.httpClient.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}

		// Wrap with request signing if enabled
		if client.requestSigner != nil {
			transport = &signingTransport{
				transport: transport,
				signer:    client.requestSigner,
			}
		}

		// Wrap with rate limiting if enabled
		if client.rateLimiter != nil {
			transport = &rateLimitTransport{
				transport:   transport,
				rateLimiter: client.rateLimiter,
			}
		}

		client.httpClient.Transport = transport
	}

	return client
}

// makeRequest performs an HTTP request with authentication
func (c *Client) makeRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	// Get valid JWT token
	token, err := c.authManager.GetToken(ctx)
	if err != nil {
		return &errors.AuthenticationError{
			Message: "Failed to obtain authentication token",
			Err:     err,
		}
	}

	// Prepare request body
	var reqBody io.Reader
	if body != nil {
		// Validate request body if it has validation tags
		if err := c.validator.Struct(body); err != nil {
			return &errors.ValidationError{
				Field:   "request",
				Message: err.Error(),
			}
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			return pkgerrors.Wrap(err, "failed to marshal request body")
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create request
	fullURL := c.config.BaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to create request")
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &errors.NetworkError{
			Operation: fmt.Sprintf("%s %s", method, endpoint),
			Err:       err,
		}
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to read response body")
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		var errorResp models.ErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err == nil {
			return errors.NewAPIError(errorResp.StatusCode, errorResp.Message, errorResp.Details, endpoint)
		}
		return errors.NewAPIError(resp.StatusCode, "Unknown error", string(respBody), endpoint)
	}

	// Parse successful response
	if result != nil {
		var apiResp models.APIResponse
		if err := json.Unmarshal(respBody, &apiResp); err != nil {
			return pkgerrors.Wrap(err, "failed to parse API response")
		}

		// Marshal and unmarshal data to convert to target type
		dataJSON, err := json.Marshal(apiResp.Data)
		if err != nil {
			return pkgerrors.Wrap(err, "failed to marshal response data")
		}

		if err := json.Unmarshal(dataJSON, result); err != nil {
			return pkgerrors.Wrap(err, "failed to unmarshal response data")
		}
	}

	return nil
}

// makeRequestWithQuery performs an HTTP request with query parameters
func (c *Client) makeRequestWithQuery(ctx context.Context, method, endpoint string, query url.Values, body interface{}, result interface{}) error {
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}
	return c.makeRequest(ctx, method, endpoint, body, result)
}

// GET performs a GET request
func (c *Client) GET(ctx context.Context, endpoint string, result interface{}) error {
	return c.makeRequest(ctx, "GET", endpoint, nil, result)
}

// GETWithQuery performs a GET request with query parameters
func (c *Client) GETWithQuery(ctx context.Context, endpoint string, query url.Values, result interface{}) error {
	return c.makeRequestWithQuery(ctx, "GET", endpoint, query, nil, result)
}

// POST performs a POST request
func (c *Client) POST(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.makeRequest(ctx, "POST", endpoint, body, result)
}

// PUT performs a PUT request
func (c *Client) PUT(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	return c.makeRequest(ctx, "PUT", endpoint, body, result)
}

// DELETE performs a DELETE request
func (c *Client) DELETE(ctx context.Context, endpoint string, result interface{}) error {
	return c.makeRequest(ctx, "DELETE", endpoint, nil, result)
}

// SetRetryPolicy sets a retry policy for the HTTP client
func (c *Client) SetRetryPolicy(maxRetries int, retryDelay time.Duration) {
	originalTransport := c.httpClient.Transport
	if originalTransport == nil {
		originalTransport = http.DefaultTransport
	}

	c.httpClient.Transport = &retryTransport{
		transport:  originalTransport,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// retryTransport implements automatic retry logic
type retryTransport struct {
	transport  http.RoundTripper
	maxRetries int
	retryDelay time.Duration
}

func (rt *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= rt.maxRetries; i++ {
		// Clone request body for retries
		var bodyBytes []byte
		if req.Body != nil {
			bodyBytes, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		resp, err = rt.transport.RoundTrip(req)

		// Don't retry on success or client errors (4xx)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		// Don't retry on the last attempt
		if i == rt.maxRetries {
			break
		}

		// Reset request body for retry
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Wait before retry with exponential backoff
		time.Sleep(rt.retryDelay * time.Duration(1<<uint(i)))
	}

	return resp, err
}

// SetRateLimit configures rate limiting for the HTTP client
func (c *Client) SetRateLimit(requestsPerSecond float64, burstSize int) {
	rateLimitConfig := &RateLimitConfig{
		RequestsPerSecond: requestsPerSecond,
		BurstSize:         burstSize,
		Enabled:           true,
	}
	
	c.rateLimiter = NewRateLimiter(rateLimitConfig)
	c.config.RateLimit = rateLimitConfig

	// Update HTTP client transport with rate limiting
	originalTransport := c.httpClient.Transport
	if originalTransport == nil {
		originalTransport = http.DefaultTransport
	}

	// Remove existing rate limit transport if any
	if rt, ok := originalTransport.(*rateLimitTransport); ok {
		originalTransport = rt.transport
	}
	if rt, ok := originalTransport.(*retryTransport); ok {
		if rlt, ok := rt.transport.(*rateLimitTransport); ok {
			rt.transport = rlt.transport
		}
	}

	c.httpClient.Transport = &rateLimitTransport{
		transport:   originalTransport,
		rateLimiter: c.rateLimiter,
	}
}

// EnableRateLimit enables rate limiting with current or default settings
func (c *Client) EnableRateLimit() {
	if c.config.RateLimit == nil {
		c.SetRateLimit(10.0, 20) // Default settings
	} else {
		c.config.RateLimit.Enabled = true
		c.rateLimiter = NewRateLimiter(c.config.RateLimit)

		// Update HTTP client transport
		originalTransport := c.httpClient.Transport
		if originalTransport == nil {
			originalTransport = http.DefaultTransport
		}

		c.httpClient.Transport = &rateLimitTransport{
			transport:   originalTransport,
			rateLimiter: c.rateLimiter,
		}
	}
}

// DisableRateLimit disables rate limiting
func (c *Client) DisableRateLimit() {
	if c.config.RateLimit != nil {
		c.config.RateLimit.Enabled = false
	}
	c.rateLimiter = nil

	// Remove rate limiting from HTTP client transport
	if rt, ok := c.httpClient.Transport.(*rateLimitTransport); ok {
		c.httpClient.Transport = rt.transport
	}
}

// GetRateLimiterStatus returns information about the current rate limiter state
func (c *Client) GetRateLimiterStatus() map[string]interface{} {
	status := map[string]interface{}{
		"enabled": false,
	}

	if c.config.RateLimit != nil {
		status["enabled"] = c.config.RateLimit.Enabled
		status["requestsPerSecond"] = c.config.RateLimit.RequestsPerSecond
		status["burstSize"] = c.config.RateLimit.BurstSize

		if c.rateLimiter != nil {
			status["availableTokens"] = c.rateLimiter.GetAvailableTokens()
		}
	}

	return status
}

// EnableCredentialEncryption enables credential encryption for the client
func (c *Client) EnableCredentialEncryption(encryptionPassword string) {
	if c.config.Security == nil {
		c.config.Security = &SecurityConfig{}
	}
	
	c.config.Security.EncryptCredentials = true
	c.config.Security.EncryptionPassword = encryptionPassword
	
	c.credentialManager = NewCredentialManager(
		encryptionPassword,
		c.config.Security.CredentialStore,
	)
}

// EnableRequestSigning enables request signing for the client
func (c *Client) EnableRequestSigning(signingSecret string) {
	if c.config.Security == nil {
		c.config.Security = &SecurityConfig{}
	}
	
	c.config.Security.EnableRequestSigning = true
	c.config.Security.SigningSecret = signingSecret
	
	c.requestSigner = NewRequestSigner(signingSecret)
	
	// Update transport chain
	c.updateTransportChain()
}

// DisableRequestSigning disables request signing
func (c *Client) DisableRequestSigning() {
	if c.config.Security != nil {
		c.config.Security.EnableRequestSigning = false
	}
	c.requestSigner = nil
	
	// Update transport chain
	c.updateTransportChain()
}

// updateTransportChain rebuilds the HTTP transport chain with current settings
func (c *Client) updateTransportChain() {
	transport := http.DefaultTransport
	
	// Wrap with request signing if enabled
	if c.requestSigner != nil {
		transport = &signingTransport{
			transport: transport,
			signer:    c.requestSigner,
		}
	}
	
	// Wrap with rate limiting if enabled
	if c.rateLimiter != nil {
		transport = &rateLimitTransport{
			transport:   transport,
			rateLimiter: c.rateLimiter,
		}
	}
	
	c.httpClient.Transport = transport
}

// StoreSecureCredentials encrypts and stores credentials if encryption is enabled
func (c *Client) StoreSecureCredentials(key, username, password string) error {
	if c.credentialManager == nil {
		return pkgerrors.New("credential encryption not enabled")
	}
	return c.credentialManager.StoreCredentials(key, username, password)
}

// RetrieveSecureCredentials retrieves and decrypts stored credentials
func (c *Client) RetrieveSecureCredentials(key string) (username, password string, err error) {
	if c.credentialManager == nil {
		return "", "", pkgerrors.New("credential encryption not enabled")
	}
	return c.credentialManager.RetrieveCredentials(key)
}

// GetSecurityStatus returns information about enabled security features
func (c *Client) GetSecurityStatus() map[string]interface{} {
	status := map[string]interface{}{
		"credentialEncryption": false,
		"requestSigning":       false,
		"rateLimiting":         false,
	}
	
	if c.config.Security != nil {
		status["credentialEncryption"] = c.config.Security.EncryptCredentials
		status["requestSigning"] = c.config.Security.EnableRequestSigning
	}
	
	if c.config.RateLimit != nil {
		status["rateLimiting"] = c.config.RateLimit.Enabled
	}
	
	return status
}