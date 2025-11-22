package abhi

import (
	"time"
	
	"abhi-go-sdk/client"
	"abhi-go-sdk/services"
)

// SDK represents the main Abhi SDK client
type SDK struct {
	client       *client.Client
	Employee     *services.EmployeeService
	Transaction  *services.TransactionService
	Organization *services.OrganizationService
	Misc         *services.MiscService
	Repayment    *services.RepaymentService
	Auth         *services.AuthService
}

// New creates a new Abhi SDK instance
func New(config *client.Config) *SDK {
	c := client.New(config)
	
	return &SDK{
		client:       c,
		Employee:     services.NewEmployeeService(c),
		Transaction:  services.NewTransactionService(c),
		Organization: services.NewOrganizationService(c),
		Misc:         services.NewMiscService(c),
		Repayment:    services.NewRepaymentService(c),
		Auth:         services.NewAuthService(c),
	}
}

// NewWithCredentials creates a new Abhi SDK instance with credentials
func NewWithCredentials(baseURL, username, password string) *SDK {
	config := client.NewConfig(baseURL, username, password)
	return New(config)
}

// NewForUAT creates a new SDK instance configured for UAT environment
func NewForUAT(username, password string) *SDK {
	config := client.DefaultConfig()
	config.BaseURL = "https://api-uat-v2.abhi.ae/uat-open-api"
	config.Username = username
	config.Password = password
	return New(config)
}

// NewForProduction creates a new SDK instance configured for production environment
func NewForProduction(username, password string) *SDK {
	config := client.DefaultConfig()
	config.BaseURL = "https://api.abhi.ae/open-api" // Replace with actual production URL
	config.Username = username
	config.Password = password
	return New(config)
}

// SetRetryPolicy configures retry behavior for API requests
func (s *SDK) SetRetryPolicy(maxRetries int, retryDelay int) *SDK {
	s.client.SetRetryPolicy(maxRetries, time.Duration(retryDelay)*time.Second)
	return s
}

// GetClient returns the underlying HTTP client for advanced usage
func (s *SDK) GetClient() *client.Client {
	return s.client
}

// SetRateLimit configures rate limiting for API requests
func (s *SDK) SetRateLimit(requestsPerSecond float64, burstSize int) *SDK {
	s.client.SetRateLimit(requestsPerSecond, burstSize)
	return s
}

// EnableRateLimit enables rate limiting with default settings
func (s *SDK) EnableRateLimit() *SDK {
	s.client.EnableRateLimit()
	return s
}

// DisableRateLimit disables rate limiting
func (s *SDK) DisableRateLimit() *SDK {
	s.client.DisableRateLimit()
	return s
}

// GetRateLimiterStatus returns information about the current rate limiter state
func (s *SDK) GetRateLimiterStatus() map[string]interface{} {
	return s.client.GetRateLimiterStatus()
}

// EnableCredentialEncryption enables credential encryption with the given password
func (s *SDK) EnableCredentialEncryption(encryptionPassword string) *SDK {
	s.client.EnableCredentialEncryption(encryptionPassword)
	return s
}

// EnableRequestSigning enables request signing for additional security
func (s *SDK) EnableRequestSigning(signingSecret string) *SDK {
	s.client.EnableRequestSigning(signingSecret)
	return s
}

// DisableRequestSigning disables request signing
func (s *SDK) DisableRequestSigning() *SDK {
	s.client.DisableRequestSigning()
	return s
}

// StoreSecureCredentials encrypts and stores credentials
func (s *SDK) StoreSecureCredentials(key, username, password string) error {
	return s.client.StoreSecureCredentials(key, username, password)
}

// RetrieveSecureCredentials retrieves and decrypts stored credentials
func (s *SDK) RetrieveSecureCredentials(key string) (username, password string, err error) {
	return s.client.RetrieveSecureCredentials(key)
}

// GetSecurityStatus returns information about enabled security features
func (s *SDK) GetSecurityStatus() map[string]interface{} {
	return s.client.GetSecurityStatus()
}