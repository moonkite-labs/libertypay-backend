package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"abhi-go-sdk"
	"abhi-go-sdk/errors"
	"abhi-go-sdk/models"
)

func main() {
	// Initialize the SDK
	sdk := abhi.NewForUAT("your-username", "your-password")

	// Set retry policy (optional)
	sdk.SetRetryPolicy(3, 2) // 3 retries with 2 second delay

	// Create context
	ctx := context.Background()

	// Employee Management Examples
	fmt.Println("=== Employee Management Examples ===")
	
	// Create a new employee
	if err := createEmployeeExample(ctx, sdk); err != nil {
		log.Printf("Error creating employee: %v", err)
	}

	// List employees
	if err := listEmployeesExample(ctx, sdk); err != nil {
		log.Printf("Error listing employees: %v", err)
	}

	// Search employees
	if err := searchEmployeeExample(ctx, sdk); err != nil {
		log.Printf("Error searching employees: %v", err)
	}

	// Transaction Management Examples
	fmt.Println("\n=== Transaction Management Examples ===")

	// Validate transaction
	if err := validateTransactionExample(ctx, sdk); err != nil {
		log.Printf("Error validating transaction: %v", err)
	}

	// Create advance transaction
	if err := createAdvanceTransactionExample(ctx, sdk); err != nil {
		log.Printf("Error creating advance transaction: %v", err)
	}

	// Get employee transaction history
	if err := getTransactionHistoryExample(ctx, sdk); err != nil {
		log.Printf("Error getting transaction history: %v", err)
	}

	// Get monthly balance
	if err := getMonthlyBalanceExample(ctx, sdk); err != nil {
		log.Printf("Error getting monthly balance: %v", err)
	}

	// Employer transaction management
	if err := employerTransactionExample(ctx, sdk); err != nil {
		log.Printf("Error with employer transactions: %v", err)
	}
}

// createEmployeeExample demonstrates creating a new employee
func createEmployeeExample(ctx context.Context, sdk *abhi.SDK) error {
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

	err := sdk.Employee.CreateSingle(ctx, employee)
	if err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}

	fmt.Printf("✓ Employee created successfully: %s %s\n", employee.FirstName, employee.LastName)
	return nil
}

// listEmployeesExample demonstrates listing employees with pagination
func listEmployeesExample(ctx context.Context, sdk *abhi.SDK) error {
	opts := &models.EmployeeListOptions{
		Page:  1,
		Limit: 10,
	}

	result, err := sdk.Employee.List(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list employees: %w", err)
	}

	fmt.Printf("✓ Found %d employees (showing %d)\n", result.Total, len(result.Results))
	for _, emp := range result.Results {
		fmt.Printf("  - %s %s (%s) - %s\n", emp.FirstName, emp.LastName, emp.EmployeeCode, emp.Department)
	}

	return nil
}

// searchEmployeeExample demonstrates searching for employees
func searchEmployeeExample(ctx context.Context, sdk *abhi.SDK) error {
	employees, err := sdk.Employee.Search(ctx, "Engineering", 10)
	if err != nil {
		return fmt.Errorf("failed to search employees: %w", err)
	}

	fmt.Printf("✓ Found %d employees in Engineering department\n", len(employees))
	for _, emp := range employees {
		fmt.Printf("  - %s %s - %s\n", emp.FirstName, emp.LastName, emp.Designation)
	}

	return nil
}

// validateTransactionExample demonstrates transaction validation
func validateTransactionExample(ctx context.Context, sdk *abhi.SDK) error {
	// Assuming we have an employee ID
	employeeID := "some-employee-id"
	amount := 1000.0

	req := models.TransactionValidationRequest{
		EmployeeID: employeeID,
		Amount:     amount,
	}

	result, err := sdk.Transaction.ValidateEmployeeTransaction(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to validate transaction: %w", err)
	}

	if result.IsValid {
		fmt.Printf("✓ Transaction is valid. Max amount: %.2f, Available: %.2f\n", 
			result.MaxAmount, result.AvailableAmount)
	} else {
		fmt.Printf("✗ Transaction is invalid: %s\n", result.Message)
	}

	return nil
}

// createAdvanceTransactionExample demonstrates creating an advance transaction
func createAdvanceTransactionExample(ctx context.Context, sdk *abhi.SDK) error {
	// Assuming we have an employee ID
	employeeID := "some-employee-id"
	amount := 500.0
	description := "Medical emergency advance"

	transaction, err := sdk.Transaction.CreateAdvanceTransaction(ctx, employeeID, amount, description)
	if err != nil {
		return fmt.Errorf("failed to create advance transaction: %w", err)
	}

	fmt.Printf("✓ Advance transaction created: ID=%s, Amount=%.2f, Status=%s\n",
		transaction.ID, transaction.Amount, transaction.Status)

	return nil
}

// getTransactionHistoryExample demonstrates getting employee transaction history
func getTransactionHistoryExample(ctx context.Context, sdk *abhi.SDK) error {
	// Assuming we have an employee ID
	employeeID := "some-employee-id"

	opts := &models.TransactionListOptions{
		Page:  1,
		Limit: 10,
		Type:  "advance", // Only advance transactions
	}

	history, err := sdk.Transaction.GetEmployeeTransactionHistory(ctx, employeeID, opts)
	if err != nil {
		return fmt.Errorf("failed to get transaction history: %w", err)
	}

	fmt.Printf("✓ Found %d transactions for employee %s\n", history.TotalCount, employeeID)
	for _, tx := range history.Transactions {
		fmt.Printf("  - %s: %.2f AED (%s) - %s\n", 
			tx.Type, tx.Amount, tx.Status, tx.RequestedAt.Format("2006-01-02"))
	}

	return nil
}

// getMonthlyBalanceExample demonstrates getting employee monthly balance
func getMonthlyBalanceExample(ctx context.Context, sdk *abhi.SDK) error {
	// Assuming we have an employee ID
	employeeID := "some-employee-id"

	// Get current month/year
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	balance, err := sdk.Transaction.GetEmployeeMonthlyBalance(ctx, employeeID, month, year)
	if err != nil {
		return fmt.Errorf("failed to get monthly balance: %w", err)
	}

	fmt.Printf("✓ Monthly balance for %s (%d/%d):\n", employeeID, month, year)
	fmt.Printf("  - Net Salary: %.2f AED\n", balance.Balance.NetSalary)
	fmt.Printf("  - Available: %.2f AED\n", balance.Balance.AvailableAmount)
	fmt.Printf("  - Used: %.2f AED\n", balance.Balance.UsedAmount)
	fmt.Printf("  - Pending: %.2f AED\n", balance.Balance.PendingAmount)

	return nil
}

// employerTransactionExample demonstrates employer transaction management
func employerTransactionExample(ctx context.Context, sdk *abhi.SDK) error {
	// Get pending transactions
	pending, err := sdk.Transaction.GetPendingTransactions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending transactions: %w", err)
	}

	fmt.Printf("✓ Found %d pending transactions\n", len(pending))
	for _, tx := range pending {
		fmt.Printf("  - %s %s: %.2f AED (%s)\n", 
			tx.EmployeeName, tx.EmployeeCode, tx.Amount, tx.RequestedAt)
	}

	// Get transactions for a date range
	startDate := "2024-01-01"
	endDate := "2024-12-31"

	dateRangeTransactions, err := sdk.Transaction.GetTransactionsByDateRange(ctx, startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to get transactions by date range: %w", err)
	}

	fmt.Printf("✓ Found %d transactions between %s and %s\n", 
		len(dateRangeTransactions), startDate, endDate)

	return nil
}

// errorHandlingExample demonstrates error handling
func errorHandlingExample(ctx context.Context, sdk *abhi.SDK) {
	// This will likely fail and demonstrate error handling
	_, err := sdk.Employee.GetByID(ctx, "non-existent-id")
	if err != nil {
		// Type assertion to get specific error information
		if apiErr, ok := err.(*errors.APIError); ok {
			switch {
			case apiErr.IsNotFound():
				fmt.Printf("Employee not found: %s\n", apiErr.Message)
			case apiErr.IsUnauthorized():
				fmt.Printf("Authentication failed: %s\n", apiErr.Message)
			case apiErr.IsRateLimited():
				fmt.Printf("Rate limit exceeded: %s\n", apiErr.Message)
			default:
				fmt.Printf("API error: %s\n", apiErr.Error())
			}
		} else {
			fmt.Printf("Unknown error: %v\n", err)
		}
	}
}