package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"abhi-go-sdk/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

// AuthManager handles JWT token management
type AuthManager struct {
	config       *Config
	token        string
	expiresAt    time.Time
	mutex        sync.RWMutex
	httpClient   *http.Client
	refreshing   bool
	refreshMutex sync.Mutex
}

// NewAuthManager creates a new authentication manager
func NewAuthManager(config *Config) *AuthManager {
	return &AuthManager{
		config:     config,
		httpClient: config.HTTPClient,
	}
}

// GetToken returns a valid JWT token, refreshing if necessary
func (a *AuthManager) GetToken(ctx context.Context) (string, error) {
	a.mutex.RLock()
	if a.isTokenValid() {
		token := a.token
		a.mutex.RUnlock()
		return token, nil
	}
	a.mutex.RUnlock()

	return a.refreshToken(ctx)
}

// isTokenValid checks if the current token is valid and not expired
func (a *AuthManager) isTokenValid() bool {
	if a.token == "" {
		return false
	}

	// Check if token expires in the next 5 minutes (buffer time)
	return time.Now().Add(5 * time.Minute).Before(a.expiresAt)
}

// refreshToken obtains a new JWT token
func (a *AuthManager) refreshToken(ctx context.Context) (string, error) {
	a.refreshMutex.Lock()
	defer a.refreshMutex.Unlock()

	// Double-check if another goroutine already refreshed the token
	a.mutex.RLock()
	if a.isTokenValid() {
		token := a.token
		a.mutex.RUnlock()
		return token, nil
	}
	a.mutex.RUnlock()

	// Perform login to get new token
	loginReq := models.LoginRequest{
		Username: a.config.Username,
		Password: a.config.Password,
	}

	reqBody, err := json.Marshal(loginReq)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal login request")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.config.BaseURL+"/auth/login", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", errors.Wrap(err, "failed to create login request")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to perform login request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp models.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return "", fmt.Errorf("login failed: %s", errorResp.Message)
		}
		return "", fmt.Errorf("login failed with status code: %d", resp.StatusCode)
	}

	var apiResp models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", errors.Wrap(err, "failed to decode login response")
	}

	loginData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return "", errors.New("invalid login response data format")
	}

	token, ok := loginData["token"].(string)
	if !ok {
		return "", errors.New("token not found in login response")
	}

	// Parse JWT to get expiration time
	expiresAt, err := a.parseTokenExpiration(token)
	if err != nil {
		// If we can't parse expiration, set it to 23 hours from now (1 hour buffer)
		expiresAt = time.Now().Add(23 * time.Hour)
	}

	a.mutex.Lock()
	a.token = token
	a.expiresAt = expiresAt
	a.mutex.Unlock()

	return token, nil
}

// parseTokenExpiration extracts the expiration time from JWT token
func (a *AuthManager) parseTokenExpiration(tokenString string) (time.Time, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return time.Time{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return time.Time{}, errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return time.Time{}, errors.New("expiration claim not found")
	}

	return time.Unix(int64(exp), 0), nil
}

// ClearToken clears the stored token (useful for logout)
func (a *AuthManager) ClearToken() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.token = ""
	a.expiresAt = time.Time{}
}