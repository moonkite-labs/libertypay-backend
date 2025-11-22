package abhi

import (
	"testing"

	"abhi-go-sdk/client"
)

func TestNew(t *testing.T) {
	config := &client.Config{
		BaseURL:  "https://api-test.example.com",
		Username: "test-user",
		Password: "test-pass",
	}

	sdk := New(config)

	if sdk == nil {
		t.Fatal("Expected SDK to be non-nil")
	}

	if sdk.Employee == nil {
		t.Error("Expected Employee service to be initialized")
	}
	if sdk.Transaction == nil {
		t.Error("Expected Transaction service to be initialized")
	}
	if sdk.Organization == nil {
		t.Error("Expected Organization service to be initialized")
	}
	if sdk.Misc == nil {
		t.Error("Expected Misc service to be initialized")
	}
	if sdk.Repayment == nil {
		t.Error("Expected Repayment service to be initialized")
	}
	if sdk.Auth == nil {
		t.Error("Expected Auth service to be initialized")
	}
}

func TestNewWithCredentials(t *testing.T) {
	baseURL := "https://api-test.example.com"
	username := "test-user"
	password := "test-pass"

	sdk := NewWithCredentials(baseURL, username, password)

	if sdk == nil {
		t.Fatal("Expected SDK to be non-nil")
	}

	// Verify all services are initialized
	if sdk.Employee == nil {
		t.Error("Expected Employee service to be initialized")
	}
}

func TestNewForUAT(t *testing.T) {
	username := "test-user"
	password := "test-pass"

	sdk := NewForUAT(username, password)

	if sdk == nil {
		t.Fatal("Expected SDK to be non-nil")
	}

	// Verify UAT URL is set correctly
	client := sdk.GetClient()
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
}

func TestNewForProduction(t *testing.T) {
	username := "test-user"
	password := "test-pass"

	sdk := NewForProduction(username, password)

	if sdk == nil {
		t.Fatal("Expected SDK to be non-nil")
	}

	client := sdk.GetClient()
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
}

func TestSetRetryPolicy(t *testing.T) {
	sdk := NewForUAT("test-user", "test-pass")
	
	// Test method chaining
	result := sdk.SetRetryPolicy(3, 2)
	if result != sdk {
		t.Error("Expected SetRetryPolicy to return the same SDK instance for method chaining")
	}
}

func TestGetClient(t *testing.T) {
	sdk := NewForUAT("test-user", "test-pass")
	
	client := sdk.GetClient()
	if client == nil {
		t.Error("Expected GetClient to return non-nil client")
	}
}

// Benchmark tests
func BenchmarkNewForUAT(b *testing.B) {
	username := "test-user"
	password := "test-pass"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sdk := NewForUAT(username, password)
		_ = sdk
	}
}

func BenchmarkSetRetryPolicy(b *testing.B) {
	sdk := NewForUAT("test-user", "test-pass")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sdk.SetRetryPolicy(3, 2)
	}
}