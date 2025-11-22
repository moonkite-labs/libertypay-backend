// Example demonstrating security features
// Run with: go run examples/security_demo.go
package security_demo

import (
	"context"
	"fmt"
	"log"

	"abhi-go-sdk"
	"abhi-go-sdk/client"
)

func mainSecurityDemo() {
	// Example 1: Basic SDK with security features
	fmt.Println("=== Security Features Example ===")

	// Initialize SDK
	sdk := abhi.NewForUAT("your-username", "your-password")

	// Enable rate limiting (10 requests per second, burst of 20)
	sdk.SetRateLimit(10.0, 20)
	fmt.Println("✓ Rate limiting enabled")

	// Enable credential encryption
	sdk.EnableCredentialEncryption("your-strong-encryption-password")
	fmt.Println("✓ Credential encryption enabled")

	// Enable request signing for additional security
	sdk.EnableRequestSigning("your-secret-signing-key")
	fmt.Println("✓ Request signing enabled")

	// Store encrypted credentials
	err := sdk.StoreSecureCredentials("production", "prod-username", "prod-password")
	if err != nil {
		log.Printf("Failed to store credentials: %v", err)
	} else {
		fmt.Println("✓ Credentials stored securely")
	}

	// Retrieve encrypted credentials
	username, _, err := sdk.RetrieveSecureCredentials("production")
	if err != nil {
		log.Printf("Failed to retrieve credentials: %v", err)
	} else {
		fmt.Printf("✓ Retrieved credentials: %s / %s\n", username, "***")
	}

	// Check security status
	status := sdk.GetSecurityStatus()
	fmt.Printf("✓ Security Status: %+v\n", status)

	// Check rate limiter status
	rateLimiterStatus := sdk.GetRateLimiterStatus()
	fmt.Printf("✓ Rate Limiter Status: %+v\n", rateLimiterStatus)

	// Example 2: Advanced configuration
	fmt.Println("\n=== Advanced Security Configuration ===")

	config := client.DefaultConfig()
	config.BaseURL = "https://api-uat-v2.abhi.ae/uat-open-api"
	config.Username = "test-user"
	config.Password = "test-pass"

	// Configure security settings
	config.EnableCredentialEncryption("encryption-password")
	config.EnableRequestSigning("signing-secret")
	config.SetRateLimit(15.0, 30)

	// Create SDK with advanced config
	advancedSDK := abhi.New(config)
	fmt.Println("✓ Advanced SDK created with security features")

	// Use the SDK for API calls (example with context)
	ctx := context.Background()
	
	// Example API call with all security features enabled
	fmt.Println("\n=== Making Secure API Calls ===")
	
	// List employees with all security features
	_, err = advancedSDK.Employee.List(ctx, nil)
	if err != nil {
		fmt.Printf("Note: API call failed (expected with test credentials): %v\n", err)
	} else {
		fmt.Printf("✓ Successfully retrieved employees\n")
	}

	// Example 3: Dynamic security control
	fmt.Println("\n=== Dynamic Security Control ===")

	// Start with basic SDK
	dynamicSDK := abhi.NewForUAT("user", "pass")

	// Enable features dynamically
	dynamicSDK.EnableRateLimit()
	fmt.Println("✓ Rate limiting enabled dynamically")

	// Enable request signing
	dynamicSDK.EnableRequestSigning("dynamic-secret")
	fmt.Println("✓ Request signing enabled dynamically")

	// Disable request signing
	dynamicSDK.DisableRequestSigning()
	fmt.Println("✓ Request signing disabled")

	// Disable rate limiting
	dynamicSDK.DisableRateLimit()
	fmt.Println("✓ Rate limiting disabled")

	// Final status check
	finalStatus := dynamicSDK.GetSecurityStatus()
	fmt.Printf("✓ Final Security Status: %+v\n", finalStatus)
}

// Demo function - call this from main package to run demos
func RunSecurityDemo() {
	mainSecurityDemo()
	demonstrateCredentialEncryption()
	demonstrateRateLimiting()
}

// demonstrateCredentialEncryption shows credential encryption features
func demonstrateCredentialEncryption() {
	fmt.Println("\n=== Credential Encryption Demo ===")

	// Create credential manager
	credManager := client.NewCredentialManager("strong-password", nil)

	// Store credentials
	err := credManager.StoreCredentials("api-key", "my-username", "my-password")
	if err != nil {
		log.Printf("Error storing credentials: %v", err)
		return
	}
	fmt.Println("✓ Credentials encrypted and stored")

	// Retrieve credentials
	username, _, err := credManager.RetrieveCredentials("api-key")
	if err != nil {
		log.Printf("Error retrieving credentials: %v", err)
		return
	}
	fmt.Printf("✓ Retrieved: %s / %s\n", username, "***")

	// Check if credentials exist
	exists := credManager.CredentialsExist("api-key")
	fmt.Printf("✓ Credentials exist: %v\n", exists)

	// Delete credentials
	err = credManager.DeleteCredentials("api-key")
	if err != nil {
		log.Printf("Error deleting credentials: %v", err)
		return
	}
	fmt.Println("✓ Credentials deleted")
}

// demonstrateRateLimiting shows rate limiting features
func demonstrateRateLimiting() {
	fmt.Println("\n=== Rate Limiting Demo ===")

	// Create rate limiter
	rateLimitConfig := &client.RateLimitConfig{
		RequestsPerSecond: 5.0, // 5 requests per second
		BurstSize:         10,  // burst of 10 requests
		Enabled:           true,
	}
	
	rateLimiter := client.NewRateLimiter(rateLimitConfig)

	// Check available tokens
	tokens := rateLimiter.GetAvailableTokens()
	fmt.Printf("✓ Available tokens: %.2f\n", tokens)

	// Test rate limiting
	for i := 0; i < 15; i++ {
		if rateLimiter.Allow() {
			fmt.Printf("✓ Request %d allowed\n", i+1)
		} else {
			fmt.Printf("✗ Request %d rate limited\n", i+1)
		}
	}
}