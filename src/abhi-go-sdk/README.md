# Abhi Go SDK

A comprehensive Go SDK for integrating with the Abhi Open API for Early Wage Access (EWA) services in the UAE.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](#testing)

## 🚀 Features

### Core Functionality
- 🔐 **Automatic JWT Authentication** - Handles token management and refresh
- 👥 **Employee Management** - Complete CRUD operations for employees
- 💰 **Transaction Management** - Advance requests, repayments, and status tracking
- 🏢 **Organization Management** - Multi-tenant organization hierarchy support
- 🏦 **Master Data APIs** - Banks and business types management
- 💵 **Repayment Tracking** - Outstanding balances and repayment processing
- 🔑 **Multi-Login Support** - Employee, employer, and third-party authentication

### Security & Performance
- 🛡️ **Advanced Security** - Request signing with HMAC-SHA256
- 🔒 **Credential Encryption** - AES-GCM encryption for sensitive data
- ⚡ **Rate Limiting** - Token bucket algorithm with configurable limits
- 🔄 **Automatic Retry Logic** - Configurable retry policy with exponential backoff
- ✅ **Input Validation** - Built-in validation for all API requests
- 🌍 **Multi-Environment Support** - UAT and Production configurations
- 📊 **Comprehensive Error Handling** - Structured error types with detailed information

## 📦 Installation

```bash
go get github.com/your-org/abhi-go-sdk
```

## 🏁 Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "abhi-go-sdk"
    "abhi-go-sdk/models"
)

func main() {
    // Initialize SDK for UAT environment
    sdk := abhi.NewForUAT("your-username", "your-password")
    
    // Enable security features
    sdk.SetRateLimit(10.0, 20)                          // 10 req/sec, burst 20
    sdk.EnableRequestSigning("your-signing-secret")      // Request signing
    sdk.EnableCredentialEncryption("encryption-pass")   // Credential encryption
    
    ctx := context.Background()
    
    // List employees
    employees, err := sdk.Employee.List(ctx, &models.EmployeeListOptions{
        Page:  1,
        Limit: 10,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d employees\n", employees.Total)
}
```

## ⚙️ Configuration

### Environment Setup

```go
// UAT Environment (recommended for testing)
sdk := abhi.NewForUAT("username", "password")

// Production Environment
sdk := abhi.NewForProduction("username", "password")

// Custom Configuration
config := &client.Config{
    BaseURL:  "https://api-uat-v2.abhi.ae/uat-open-api",
    Username: "your-username", 
    Password: "your-password",
    Timeout:  30 * time.Second,
}
sdk := abhi.New(config)
```

### Security Configuration

```go
// Enable all security features
sdk := abhi.NewForUAT("username", "password")

// Rate Limiting: 10 requests/second, burst of 20
sdk.SetRateLimit(10.0, 20)

// Request Signing for tamper protection
sdk.EnableRequestSigning("your-secret-key")

// Credential Encryption for secure storage
sdk.EnableCredentialEncryption("encryption-password")

// Check security status
status := sdk.GetSecurityStatus()
fmt.Printf("Security Status: %+v\n", status)
```

### Advanced Configuration

```go
config := client.DefaultConfig()
config.BaseURL = "https://api-uat-v2.abhi.ae/uat-open-api"
config.Username = "your-username"
config.Password = "your-password"

// Configure rate limiting
config.SetRateLimit(15.0, 30)

// Configure security
config.EnableRequestSigning("signing-secret")
config.EnableCredentialEncryption("encryption-password")

// Custom HTTP client
config.SetHTTPClient(&http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns: 10,
    },
})

sdk := abhi.New(config)
```

## 👥 Employee Management

### Creating Employees

```go
employee := models.Employee{
    EmployeeCode:    "EMP001",
    FirstName:       "John",
    LastName:        "Doe", 
    Department:      "Engineering",
    Designation:     "Software Engineer",
    Phone:           "+971501234567",
    Email:           "john.doe@company.com",
    DOB:             "1990-01-15",
    DateOfJoining:   "2024-01-01",
    AccountTitle:    "John Doe",
    AccountNumber:   "1234567890",
    NetSalary:       "8000",
    EmiratesID:      "784-1990-1234567-1",
    Gender:          "Male",
    BankID:          "9b5fcf65-5fca-4acf-a3a5-6f79055644e1",
    PayrollStartDay: 1,
}

// Create single employee
err := sdk.Employee.CreateSingle(ctx, employee)
if err != nil {
    log.Fatal(err)
}

// Create multiple employees
employees := []models.Employee{employee1, employee2}
err := sdk.Employee.Create(ctx, employees)
```

### Employee Operations

```go
// List with pagination and filters
opts := &models.EmployeeListOptions{
    Page:       1,
    Limit:      50,
    Department: "Engineering",
    Search:     "john",
}
result, err := sdk.Employee.List(ctx, opts)

// Get by ID
employee, err := sdk.Employee.GetByID(ctx, "employee-id")

// Search employees
employees, err := sdk.Employee.Search(ctx, "software engineer", 10)

// Update employee
employee.Department = "Product"
err := sdk.Employee.UpdateSingle(ctx, employee)

// Delete employee
err := sdk.Employee.Delete(ctx, "employee-id")
```

## 💰 Transaction Management

### Creating Transactions

```go
// Create advance transaction
transaction, err := sdk.Transaction.CreateAdvanceTransaction(
    ctx, 
    "employee-id", 
    1000.0, 
    "Medical emergency"
)

// Validate transaction before creating
validation, err := sdk.Transaction.ValidateEmployeeTransaction(ctx, 
    models.TransactionValidationRequest{
        EmployeeID: "employee-id",
        Amount:     1000.0,
    })

if validation.IsValid {
    // Proceed with transaction creation
    fmt.Printf("Max amount: %.2f, Available: %.2f\n", 
        validation.MaxAmount, validation.AvailableAmount)
}
```

### Transaction History & Balance

```go
// Get employee transaction history
history, err := sdk.Transaction.GetEmployeeTransactionHistory(ctx, 
    "employee-id", 
    &models.TransactionListOptions{
        Page:  1,
        Limit: 20,
        Type:  "advance",
    })

// Get monthly balance
balance, err := sdk.Transaction.GetEmployeeMonthlyBalance(ctx, 
    "employee-id", 11, 2024)

fmt.Printf("Available: %.2f AED, Used: %.2f AED\n", 
    balance.Balance.AvailableAmount, balance.Balance.UsedAmount)
```

### Employer Transaction Management

```go
// Get pending transactions for approval
pending, err := sdk.Transaction.GetPendingTransactions(ctx)

// Get transactions by date range
transactions, err := sdk.Transaction.GetTransactionsByDateRange(ctx, 
    "2024-01-01", "2024-12-31")

// List all transactions with filters
opts := &models.EmployerTransactionListOptions{
    Page:       1,
    Limit:      50,
    Status:     "approved",
    Department: "Engineering",
}
result, err := sdk.Transaction.GetEmployerTransactions(ctx, opts)
```

## 🏢 Organization Management

### Creating Organizations

```go
orgRequest := models.CreateOrganizationRequest{
    Name:            "Tech Solutions Ltd",
    Industry:        "Technology", 
    BusinessTypeID:  "business-type-uuid",
    Address:         "123 Business Street",
    City:            "Dubai",
    ManagementAlias: "tech_solutions",
    CreditLimit:     1000000.0,
    Email:           "admin@techsolutions.com",
    Phone:           "+971501234567",
    PayrollStartDay: 1,
}

response, err := sdk.Organization.Create(ctx, orgRequest)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created organization: %s\n", response.Data.OrganizationID)
fmt.Printf("Admin user: %s\n", response.Data.Users.Admin.Username)
```

### Organization Operations

```go
// List organizations with filters
opts := &models.OrganizationListOptions{
    Page:         1,
    Limit:        20,
    ShowInactive: false,
    From:         "2024-01-01",
    To:           "2024-12-31",
}
orgs, err := sdk.Organization.List(ctx, opts)

// Get organization statistics
stats, err := sdk.Organization.GetStatistics(ctx)

// Search by industry
techOrgs, err := sdk.Organization.GetByIndustry(ctx, "Technology")
```

## 💵 Repayment Management

### Outstanding Balances

```go
// Get outstanding balances with filters
opts := &models.OutstandingBalanceListOptions{
    Page:      1,
    Limit:     50,
    Overdue:   true,
    MinAmount: 100.0,
}
outstanding, err := sdk.Repayment.GetOutstandingBalance(ctx, opts)

// Get employee-specific balance
balance, err := sdk.Repayment.GetEmployeeOutstandingBalance(ctx, "employee-id")

// Get overdue balances
overdueBalances, err := sdk.Repayment.GetOverdueBalances(ctx)
```

### Creating Repayments

```go
// Employee repayment
repaymentResponse, err := sdk.Repayment.CreateEmployeeRepayment(
    ctx,
    "employee-id",
    500.0,
    "REP-2024-001", 
    "Salary deduction repayment",
)

// Transaction-specific repayment
repaymentResponse, err := sdk.Repayment.CreateTransactionRepayment(
    ctx,
    "transaction-id",
    250.0,
    "REP-2024-002",
    "Partial transaction repayment",
)

// List repayments
repayments, err := sdk.Repayment.ListRepayments(ctx, &models.RepaymentListOptions{
    Status:    "completed",
    StartDate: "2024-01-01",
    EndDate:   "2024-12-31",
})
```

## 🔑 Authentication Management

### Multiple Login Types

```go
// Employee login (requires Emirates ID)
authResponse, err := sdk.Auth.LoginEmployee(ctx, 
    "employee@company.com", 
    "password", 
    "784-1990-1234567-1")

// Employer login
authResponse, err := sdk.Auth.LoginEmployer(ctx,
    "admin@company.com", 
    "password")

// Third-party system login
authResponse, err := sdk.Auth.LoginThirdParty(ctx,
    "api-username", 
    "api-password")
```

### Session Management

```go
// Get current user information
user, err := sdk.Auth.GetCurrentUser(ctx)

// Validate current token
isValid, err := sdk.Auth.ValidateToken(ctx)

// Change password
err = sdk.Auth.ChangePassword(ctx, models.ChangePasswordRequest{
    CurrentPassword: "old-password",
    NewPassword:     "new-password",
    ConfirmPassword: "new-password",
})

// Logout
err = sdk.Auth.LogoutCurrentSession(ctx)
```

### Multi-Factor Authentication

```go
// Setup MFA
mfaResponse, err := sdk.Auth.SetupMFA(ctx, models.MFASetupRequest{
    Method: "totp",
})

// Verify MFA code
authResponse, err := sdk.Auth.VerifyMFA(ctx, models.MFAVerificationRequest{
    Token: "mfa-token",
    Code:  "123456",
})

// Check MFA status
status, err := sdk.Auth.GetMFAStatus(ctx)
```

## 🏦 Master Data APIs

### Banks Management

```go
// Get all banks
banks, err := sdk.Misc.GetAllBanks(ctx)

// Get banks by country
uaeBanks, err := sdk.Misc.GetBanksByCountry(ctx, "UAE")

// Get active banks only
activeBanks, err := sdk.Misc.GetActiveBanks(ctx)

// Search banks
searchResults, err := sdk.Misc.SearchBanks(ctx, "Emirates", 10)
```

### Business Types Management

```go
// Get all business types
businessTypes, err := sdk.Misc.GetAllBusinessTypes(ctx)

// Get by country
uaeBusinessTypes, err := sdk.Misc.GetBusinessTypesByCountry(ctx, "UAE")

// Get active types only
activeTypes, err := sdk.Misc.GetActiveBusinessTypes(ctx)
```

## 🔒 Security Features

### Credential Encryption

```go
// Enable credential encryption
sdk.EnableCredentialEncryption("strong-encryption-password")

// Store encrypted credentials
err := sdk.StoreSecureCredentials("production", "prod-user", "prod-pass")

// Retrieve encrypted credentials
username, password, err := sdk.RetrieveSecureCredentials("production")

// Check if credentials exist
exists := sdk.client.credentialManager.CredentialsExist("production")
```

### Request Signing

```go
// Enable request signing for tamper protection
sdk.EnableRequestSigning("your-secret-signing-key")

// All subsequent requests will be automatically signed
// Signature includes: method, path, headers, body hash, timestamp

// Disable request signing
sdk.DisableRequestSigning()
```

### Rate Limiting

```go
// Configure rate limiting
sdk.SetRateLimit(10.0, 20) // 10 requests/second, burst of 20

// Enable with defaults (10 req/sec, burst 20)
sdk.EnableRateLimit()

// Check rate limiter status
status := sdk.GetRateLimiterStatus()
fmt.Printf("Available tokens: %.2f\n", status["availableTokens"])

// Disable rate limiting
sdk.DisableRateLimit()
```

## ⚠️ Error Handling

The SDK provides structured error handling with specific error types:

```go
employee, err := sdk.Employee.GetByID(ctx, "non-existent-id")
if err != nil {
    if apiErr, ok := err.(*errors.APIError); ok {
        switch {
        case apiErr.IsNotFound():
            fmt.Println("Employee not found")
        case apiErr.IsUnauthorized():
            fmt.Println("Authentication failed")
        case apiErr.IsRateLimited():
            fmt.Println("Rate limit exceeded")
        case apiErr.IsServerError():
            fmt.Printf("Server error: %s", apiErr.Message)
        default:
            fmt.Printf("API error [%d]: %s", apiErr.StatusCode, apiErr.Message)
        }
    }
}
```

### Error Types

- **`APIError`** - HTTP API errors with status codes and helper methods
- **`ValidationError`** - Request validation errors with field details
- **`NetworkError`** - Network connectivity issues
- **`AuthenticationError`** - Authentication/authorization errors

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./client
go test ./errors

# Run tests with verbose output
go test -v ./...
```

### Test Coverage

The SDK includes comprehensive tests covering:
- ✅ All HTTP client functionality
- ✅ Authentication and token management
- ✅ Error handling scenarios
- ✅ Rate limiting behavior
- ✅ Security features (encryption, signing)
- ✅ Configuration management

## 📊 Performance & Monitoring

### Rate Limiting Monitoring

```go
// Get current rate limiter status
status := sdk.GetRateLimiterStatus()
fmt.Printf("Rate limiting enabled: %v\n", status["enabled"])
fmt.Printf("Requests per second: %v\n", status["requestsPerSecond"])
fmt.Printf("Available tokens: %v\n", status["availableTokens"])
```

### Security Status

```go
// Get security feature status
securityStatus := sdk.GetSecurityStatus()
fmt.Printf("Credential encryption: %v\n", securityStatus["credentialEncryption"])
fmt.Printf("Request signing: %v\n", securityStatus["requestSigning"])
fmt.Printf("Rate limiting: %v\n", securityStatus["rateLimiting"])
```

### Retry Configuration

```go
// Configure retry policy: 3 retries with 2 second initial delay
sdk.SetRetryPolicy(3, 2)

// Automatic exponential backoff:
// - 1st retry: 2 seconds
// - 2nd retry: 4 seconds  
// - 3rd retry: 8 seconds
```

## 🌍 Environment Support

| Environment | URL | Description |
|------------|-----|-------------|
| UAT | `https://api-uat-v2.abhi.ae/uat-open-api` | Testing environment |
| Production | `https://api.abhi.ae/open-api` | Live environment |

```go
// Environment-specific initialization
sdk := abhi.NewForUAT("username", "password")        // UAT
sdk := abhi.NewForProduction("username", "password")  // Production
```

## 📚 Examples

Complete examples are available in the `examples/` directory:

- **`examples/basic/main.go`** - Basic usage examples
- **`examples/advanced/extended.go`** - Advanced features and error handling
- **`examples/security/security_demo.go`** - Security features demonstration

Run examples:

```bash
# Basic example
go run examples/basic/main.go

# Advanced features
go run examples/advanced/extended.go

# Security features demo
go run examples/security/security_demo.go
```

## 🔧 Development

### Building

```bash
# Build all packages
go build ./...

# Build specific package
go build ./client

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build ./...
```

### Dependencies

The SDK uses minimal, well-maintained dependencies:

- `github.com/golang-jwt/jwt/v4` - JWT token handling
- `github.com/go-playground/validator/v10` - Input validation
- `github.com/pkg/errors` - Enhanced error handling

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Add tests for your changes
4. Ensure all tests pass (`go test ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Contribution Guidelines

- Follow Go conventions and best practices
- Add unit tests for new functionality
- Update documentation and examples
- Ensure backward compatibility
- Write clear, descriptive commit messages

## 📋 API Reference

### SDK Methods

| Method | Description |
|--------|-------------|
| `New(config)` | Create SDK with custom config |
| `NewForUAT(username, password)` | Create SDK for UAT environment |
| `NewForProduction(username, password)` | Create SDK for production |
| `SetRetryPolicy(retries, delay)` | Configure retry behavior |
| `SetRateLimit(rps, burst)` | Configure rate limiting |
| `EnableRequestSigning(secret)` | Enable request signing |
| `EnableCredentialEncryption(password)` | Enable credential encryption |
| `GetSecurityStatus()` | Get security feature status |

### Service Methods

Each service provides comprehensive CRUD operations with both detailed and convenience methods:

- **Employee Service** - 15+ methods for employee management
- **Transaction Service** - 20+ methods for transaction handling  
- **Organization Service** - 12+ methods for organization management
- **Repayment Service** - 10+ methods for repayment processing
- **Auth Service** - 15+ methods for authentication
- **Misc Service** - 10+ methods for master data

## 🔐 Security Notes

- **Credentials**: Automatically managed through JWT tokens
- **Token Refresh**: Automatic refresh with 5-minute buffer
- **Request Signing**: HMAC-SHA256 with timestamp validation
- **Encryption**: AES-GCM for credential storage
- **Rate Limiting**: Token bucket algorithm prevents abuse
- **Validation**: Comprehensive input validation
- **HTTPS**: All communications over secure channels

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💬 Support

For support and questions:

- 📧 Create an issue on GitHub
- 📖 Review the API documentation
- 💡 Check the examples directory for usage patterns
- 🔍 Search existing issues for solutions

---

## 🏆 Features Comparison

| Feature | Basic SDK | Enhanced SDK (This Implementation) |
|---------|-----------|-----------------------------------|
| HTTP Client | ✅ | ✅ |
| Authentication | ✅ | ✅ Enhanced with auto-refresh |
| Error Handling | ✅ | ✅ Structured with helpers |
| Rate Limiting | ❌ | ✅ Token bucket algorithm |
| Request Signing | ❌ | ✅ HMAC-SHA256 |
| Credential Encryption | ❌ | ✅ AES-GCM |
| Retry Logic | ✅ | ✅ Enhanced with exponential backoff |
| Validation | ✅ | ✅ Comprehensive |
| Testing | ❌ | ✅ 30+ test functions |
| Security Monitoring | ❌ | ✅ Status reporting |
| Examples | ✅ Basic | ✅ Comprehensive |
| Documentation | ✅ | ✅ Extensive |

---

Made with ❤️ for the Abhi developer community