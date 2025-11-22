package main

import (
	"context"
	"fmt"
	"log"

	"abhi-go-sdk"
	"abhi-go-sdk/models"
)

func main() {
	// Initialize the SDK
	sdk := abhi.NewForUAT("your-username", "your-password")
	ctx := context.Background()

	fmt.Println("=== Extended Abhi SDK Examples ===")

	// Organization Management Examples
	if err := organizationExamples(ctx, sdk); err != nil {
		log.Printf("Organization examples error: %v", err)
	}

	// Miscellaneous API Examples  
	if err := miscExamples(ctx, sdk); err != nil {
		log.Printf("Miscellaneous examples error: %v", err)
	}

	// Repayment Examples
	if err := repaymentExamples(ctx, sdk); err != nil {
		log.Printf("Repayment examples error: %v", err)
	}

	// Authentication Examples
	if err := authExamples(ctx, sdk); err != nil {
		log.Printf("Authentication examples error: %v", err)
	}
}

// organizationExamples demonstrates organization management
func organizationExamples(ctx context.Context, sdk *abhi.SDK) error {
	fmt.Println("\n=== Organization Management ===")

	// List all organizations
	orgs, err := sdk.Organization.List(ctx, &models.OrganizationListOptions{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		return err
	}
	fmt.Printf("✓ Found %d organizations\n", orgs.Total)

	// Get active organizations
	activeOrgs, err := sdk.Organization.GetActive(ctx, nil)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Found %d active organizations\n", len(activeOrgs.Results))

	// Create a new organization
	newOrg := models.CreateOrganizationRequest{
		Name:            "Tech Innovators Ltd",
		Industry:        "Technology",
		BusinessTypeID:  "ab4d503f-bf9f-4a95-aeb3-20817776bc7c", // Example ID
		Address:         "123 Innovation Street",
		City:            "Dubai",
		ManagementAlias: "tech_innovators",
		CreditLimit:     1000000.0,
		Email:           "admin@techinnovators.com",
		Phone:           "+971501234567",
		PayrollStartDay: 1,
	}

	createResp, err := sdk.Organization.Create(ctx, newOrg)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Created organization: %s\n", createResp.Data.OrganizationID)
	fmt.Printf("  - Admin user: %s\n", createResp.Data.Users.Admin.Username)

	// Get organization statistics
	stats, err := sdk.Organization.GetStatistics(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Organization Statistics:\n")
	fmt.Printf("  - Total: %v\n", stats["total"])
	fmt.Printf("  - Active: %v\n", stats["active"])
	fmt.Printf("  - Industries: %v\n", stats["industries"])

	return nil
}

// miscExamples demonstrates miscellaneous APIs (banks, business types)
func miscExamples(ctx context.Context, sdk *abhi.SDK) error {
	fmt.Println("\n=== Banks and Business Types ===")

	// Get all banks
	banks, err := sdk.Misc.GetAllBanks(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Found %d banks\n", len(banks))

	// Get active banks only
	activeBanks, err := sdk.Misc.GetActiveBanks(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Found %d active banks\n", len(activeBanks))

	// Search banks
	emiratesBanks, err := sdk.Misc.GetBanksByCountry(ctx, "UAE")
	if err != nil {
		return err
	}
	fmt.Printf("✓ Found %d banks in UAE\n", len(emiratesBanks))

	// List first 5 banks
	for i, bank := range activeBanks {
		if i >= 5 {
			break
		}
		fmt.Printf("  - %s (%s) - %s\n", bank.Name, bank.Code, bank.Country)
	}

	// Get business types
	businessTypes, err := sdk.Misc.GetAllBusinessTypes(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Found %d business types\n", len(businessTypes))

	// Get UAE business types
	uaeBusinessTypes, err := sdk.Misc.GetBusinessTypesByCountry(ctx, "UAE")
	if err != nil {
		return err
	}
	fmt.Printf("✓ Found %d business types in UAE\n", len(uaeBusinessTypes))

	// List first 3 business types
	for i, bt := range businessTypes {
		if i >= 3 {
			break
		}
		fmt.Printf("  - %s (%s)\n", bt.Name, bt.Country)
	}

	return nil
}

// repaymentExamples demonstrates repayment management
func repaymentExamples(ctx context.Context, sdk *abhi.SDK) error {
	fmt.Println("\n=== Repayment Management ===")

	// Get outstanding balances
	outstanding, err := sdk.Repayment.GetOutstandingBalance(ctx, &models.OutstandingBalanceListOptions{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		return err
	}
	fmt.Printf("✓ Found %d employees with outstanding balances\n", outstanding.Total)

	// Show summary if available
	if outstanding.Summary.TotalEmployees > 0 {
		fmt.Printf("✓ Outstanding Balance Summary:\n")
		fmt.Printf("  - Total Employees: %d\n", outstanding.Summary.TotalEmployees)
		fmt.Printf("  - Total Outstanding: %.2f AED\n", outstanding.Summary.TotalOutstanding)
		fmt.Printf("  - Total Overdue: %.2f AED\n", outstanding.Summary.TotalOverdue)
		fmt.Printf("  - Average Outstanding: %.2f AED\n", outstanding.Summary.AverageOutstanding)
	}

	// Get overdue balances
	overdueBalances, err := sdk.Repayment.GetOverdueBalances(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Found %d employees with overdue balances\n", len(overdueBalances))

	// Show first 3 overdue balances
	for i, balance := range overdueBalances {
		if i >= 3 {
			break
		}
		fmt.Printf("  - %s: %.2f AED overdue (%d days)\n", 
			balance.EmployeeName, balance.OverdueAmount, balance.DaysPastDue)
	}

	// Create a repayment (example)
	if len(outstanding.Results) > 0 {
		employeeID := outstanding.Results[0].EmployeeID
		repaymentResp, err := sdk.Repayment.CreateEmployeeRepayment(
			ctx,
			employeeID,
			500.0,
			"REP-2024-001",
			"Partial repayment via mobile app",
		)
		if err != nil {
			// This might fail in demo, just log it
			fmt.Printf("⚠ Repayment creation failed (expected in demo): %v\n", err)
		} else {
			fmt.Printf("✓ Created repayment: %.2f AED for employee %s\n", 
				repaymentResp.Repayment.Amount, employeeID)
		}
	}

	// List recent repayments
	repayments, err := sdk.Repayment.ListRepayments(ctx, &models.RepaymentListOptions{
		Page:   1,
		Limit:  5,
		Status: "completed",
	})
	if err != nil {
		return err
	}
	fmt.Printf("✓ Found %d completed repayments\n", repayments.Total)

	return nil
}

// authExamples demonstrates authentication features  
func authExamples(ctx context.Context, sdk *abhi.SDK) error {
	fmt.Println("\n=== Authentication Examples ===")

	// Get current user info (if authenticated)
	user, err := sdk.Auth.GetCurrentUser(ctx)
	if err != nil {
		fmt.Printf("⚠ Could not get current user (expected if using API key): %v\n", err)
	} else {
		fmt.Printf("✓ Current user: %s (%s)\n", user.Username, user.Role)
		fmt.Printf("  - Organization: %s\n", user.OrganizationName)
		fmt.Printf("  - Permissions: %v\n", user.Permissions)
	}

	// Validate current token
	isValid, err := sdk.Auth.ValidateToken(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Token is valid: %t\n", isValid)

	// Example of different login types (these would fail without real credentials)
	fmt.Println("\n--- Login Examples (will fail without real credentials) ---")

	// Employee login example
	_, err = sdk.Auth.LoginEmployee(ctx, "employee@company.com", "password", "784-1990-1234567-1")
	if err != nil {
		fmt.Printf("⚠ Employee login failed (expected): %v\n", err)
	}

	// Employer login example
	_, err = sdk.Auth.LoginEmployer(ctx, "admin@company.com", "password")
	if err != nil {
		fmt.Printf("⚠ Employer login failed (expected): %v\n", err)
	}

	// Third-party login example
	_, err = sdk.Auth.LoginThirdParty(ctx, "api-user", "api-password")
	if err != nil {
		fmt.Printf("⚠ Third-party login failed (expected): %v\n", err)
	}

	// Get session info
	session, err := sdk.Auth.GetSessionInfo(ctx)
	if err != nil {
		fmt.Printf("⚠ Could not get session info: %v\n", err)
	} else {
		fmt.Printf("✓ Session info:\n")
		fmt.Printf("  - Login time: %s\n", session.LoginTime)
		fmt.Printf("  - Last activity: %s\n", session.LastActivity)
		fmt.Printf("  - Expires at: %s\n", session.ExpiresAt)
	}

	return nil
}

// demonstrateErrorHandling shows comprehensive error handling
func demonstrateErrorHandling(ctx context.Context, sdk *abhi.SDK) {
	fmt.Println("\n=== Error Handling Examples ===")

	// Try to get a non-existent organization
	_, err := sdk.Organization.GetByID(ctx, "non-existent-id")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Try to create repayment with invalid data
	invalidRepayment := models.CreateRepaymentRequest{
		Amount:                         -100, // Invalid negative amount
		ClientRepaymentReferenceNumber: "",   // Missing required field
	}

	if err := sdk.Repayment.ValidateRepayment(invalidRepayment); err != nil {
		fmt.Printf("Validation error: %v\n", err)
	}

	// Try to create organization with invalid data
	invalidOrg := models.CreateOrganizationRequest{
		Name:         "", // Missing required field
		CreditLimit:  -1000, // Invalid negative amount
	}

	if err := sdk.Organization.ValidateOrganization(invalidOrg); err != nil {
		fmt.Printf("Organization validation error: %v\n", err)
	}
}