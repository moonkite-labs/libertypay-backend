package client

import (
	"net/http"
	"time"
)

// Config holds the configuration for the Abhi API client
type Config struct {
	BaseURL           string
	Username          string
	Password          string
	HTTPClient        *http.Client
	Timeout           time.Duration
	RateLimit         *RateLimitConfig
	Security          *SecurityConfig
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	EncryptCredentials   bool
	EncryptionPassword   string
	CredentialStore      CredentialStore
	EnableRequestSigning bool
	SigningSecret        string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerSecond float64
	BurstSize         int
	Enabled           bool
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		BaseURL: "https://api-uat-v2.abhi.ae/uat-open-api",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Timeout: 30 * time.Second,
		RateLimit: &RateLimitConfig{
			RequestsPerSecond: 10.0, // Default: 10 requests per second
			BurstSize:         20,   // Default: burst of 20 requests
			Enabled:           false, // Disabled by default
		},
		Security: &SecurityConfig{
			EncryptCredentials:   false, // Disabled by default
			EnableRequestSigning: false, // Disabled by default
		},
	}
}

// NewConfig creates a new configuration with the provided base URL and credentials
func NewConfig(baseURL, username, password string) *Config {
	config := DefaultConfig()
	config.BaseURL = baseURL
	config.Username = username
	config.Password = password
	return config
}

// SetHTTPClient sets a custom HTTP client
func (c *Config) SetHTTPClient(client *http.Client) *Config {
	c.HTTPClient = client
	return c
}

// SetTimeout sets the request timeout
func (c *Config) SetTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	c.HTTPClient.Timeout = timeout
	return c
}

// SetRateLimit sets the rate limiting configuration
func (c *Config) SetRateLimit(requestsPerSecond float64, burstSize int) *Config {
	c.RateLimit = &RateLimitConfig{
		RequestsPerSecond: requestsPerSecond,
		BurstSize:         burstSize,
		Enabled:           true,
	}
	return c
}

// EnableRateLimit enables rate limiting with current settings
func (c *Config) EnableRateLimit() *Config {
	if c.RateLimit == nil {
		c.RateLimit = &RateLimitConfig{
			RequestsPerSecond: 10.0,
			BurstSize:         20,
			Enabled:           true,
		}
	} else {
		c.RateLimit.Enabled = true
	}
	return c
}

// DisableRateLimit disables rate limiting
func (c *Config) DisableRateLimit() *Config {
	if c.RateLimit != nil {
		c.RateLimit.Enabled = false
	}
	return c
}

// EnableCredentialEncryption enables credential encryption
func (c *Config) EnableCredentialEncryption(encryptionPassword string) *Config {
	if c.Security == nil {
		c.Security = &SecurityConfig{}
	}
	c.Security.EncryptCredentials = true
	c.Security.EncryptionPassword = encryptionPassword
	return c
}

// SetCredentialStore sets a custom credential store
func (c *Config) SetCredentialStore(store CredentialStore) *Config {
	if c.Security == nil {
		c.Security = &SecurityConfig{}
	}
	c.Security.CredentialStore = store
	return c
}

// EnableRequestSigning enables request signing for additional security
func (c *Config) EnableRequestSigning(signingSecret string) *Config {
	if c.Security == nil {
		c.Security = &SecurityConfig{}
	}
	c.Security.EnableRequestSigning = true
	c.Security.SigningSecret = signingSecret
	return c
}

// DisableCredentialEncryption disables credential encryption
func (c *Config) DisableCredentialEncryption() *Config {
	if c.Security != nil {
		c.Security.EncryptCredentials = false
	}
	return c
}

// DisableRequestSigning disables request signing
func (c *Config) DisableRequestSigning() *Config {
	if c.Security != nil {
		c.Security.EnableRequestSigning = false
	}
	return c
}