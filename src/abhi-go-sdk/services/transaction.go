package services

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"abhi-go-sdk/client"
	"abhi-go-sdk/models"
)

// TransactionService handles transaction-related API operations
type TransactionService struct {
	client *client.Client
}

// NewTransactionService creates a new transaction service
func NewTransactionService(client *client.Client) *TransactionService {
	return &TransactionService{
		client: client,
	}
}

// Employee Transaction Methods

// CreateEmployeeTransaction creates a new transaction for an employee
func (s *TransactionService) CreateEmployeeTransaction(ctx context.Context, req models.TransactionRequest) (*models.Transaction, error) {
	var result models.Transaction
	err := s.client.POST(ctx, "/transactions/employee", req, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create employee transaction: %w", err)
	}

	return &result, nil
}

// GetEmployeeTransactionHistory retrieves transaction history for an employee
func (s *TransactionService) GetEmployeeTransactionHistory(ctx context.Context, employeeID string, opts *models.TransactionListOptions) (*models.TransactionHistoryResponse, error) {
	query := url.Values{}
	
	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Status != "" {
			query.Set("status", opts.Status)
		}
		if opts.Type != "" {
			query.Set("type", opts.Type)
		}
		if opts.StartDate != "" {
			query.Set("startDate", opts.StartDate)
		}
		if opts.EndDate != "" {
			query.Set("endDate", opts.EndDate)
		}
	}

	endpoint := fmt.Sprintf("/transactions/employee/%s/history", employeeID)
	
	var result models.TransactionHistoryResponse
	err := s.client.GETWithQuery(ctx, endpoint, query, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee transaction history: %w", err)
	}

	return &result, nil
}

// GetEmployeeMonthlyBalance retrieves monthly balance for an employee
func (s *TransactionService) GetEmployeeMonthlyBalance(ctx context.Context, employeeID string, month, year int) (*models.MonthlyBalanceResponse, error) {
	query := url.Values{}
	if month > 0 {
		query.Set("month", strconv.Itoa(month))
	}
	if year > 0 {
		query.Set("year", strconv.Itoa(year))
	}

	endpoint := fmt.Sprintf("/transactions/employee/%s/balance", employeeID)
	
	var result models.MonthlyBalanceResponse
	err := s.client.GETWithQuery(ctx, endpoint, query, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee monthly balance: %w", err)
	}

	return &result, nil
}

// ValidateEmployeeTransaction validates a transaction before processing
func (s *TransactionService) ValidateEmployeeTransaction(ctx context.Context, req models.TransactionValidationRequest) (*models.TransactionValidationResponse, error) {
	var result models.TransactionValidationResponse
	err := s.client.POST(ctx, "/transactions/employee/validate", req, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to validate employee transaction: %w", err)
	}

	return &result, nil
}

// GetEmployeeTransactionStatus retrieves the status of a specific transaction
func (s *TransactionService) GetEmployeeTransactionStatus(ctx context.Context, transactionID string) (*models.TransactionStatusResponse, error) {
	endpoint := fmt.Sprintf("/transactions/employee/%s/status", transactionID)
	
	var result models.TransactionStatusResponse
	err := s.client.GET(ctx, endpoint, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee transaction status: %w", err)
	}

	return &result, nil
}

// Employer Transaction Methods

// GetEmployerTransactions retrieves transactions from employer perspective
func (s *TransactionService) GetEmployerTransactions(ctx context.Context, opts *models.EmployerTransactionListOptions) (*models.EmployerTransactionResponse, error) {
	query := url.Values{}
	
	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Status != "" {
			query.Set("status", opts.Status)
		}
		if opts.Type != "" {
			query.Set("type", opts.Type)
		}
		if opts.StartDate != "" {
			query.Set("startDate", opts.StartDate)
		}
		if opts.EndDate != "" {
			query.Set("endDate", opts.EndDate)
		}
		if opts.EmployeeCode != "" {
			query.Set("employeeCode", opts.EmployeeCode)
		}
		if opts.Department != "" {
			query.Set("department", opts.Department)
		}
	}

	var result models.EmployerTransactionResponse
	err := s.client.GETWithQuery(ctx, "/transactions/employer", query, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get employer transactions: %w", err)
	}

	return &result, nil
}

// GetEmployerTransactionStatus retrieves transaction status from employer perspective
func (s *TransactionService) GetEmployerTransactionStatus(ctx context.Context, transactionID string) (*models.TransactionStatusResponse, error) {
	endpoint := fmt.Sprintf("/transactions/employer/%s/status", transactionID)
	
	var result models.TransactionStatusResponse
	err := s.client.GET(ctx, endpoint, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get employer transaction status: %w", err)
	}

	return &result, nil
}

// ValidateQuestions retrieves validation questions for a transaction
func (s *TransactionService) ValidateQuestions(ctx context.Context, req models.ValidationQuestionsRequest) (*models.ValidationQuestionsResponse, error) {
	var result models.ValidationQuestionsResponse
	err := s.client.POST(ctx, "/transactions/employer/validate-questions", req, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get validation questions: %w", err)
	}

	return &result, nil
}

// SubmitValidationAnswers submits answers to validation questions
func (s *TransactionService) SubmitValidationAnswers(ctx context.Context, req models.ValidationAnswersRequest) (*models.ValidationAnswersResponse, error) {
	var result models.ValidationAnswersResponse
	err := s.client.POST(ctx, "/transactions/employer/validate-answers", req, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to submit validation answers: %w", err)
	}

	return &result, nil
}

// Convenience Methods

// GetAllEmployerTransactions retrieves all transactions with pagination handling
func (s *TransactionService) GetAllEmployerTransactions(ctx context.Context, opts *models.EmployerTransactionListOptions) ([]models.EmployerTransaction, error) {
	var allTransactions []models.EmployerTransaction
	page := 1
	limit := 100

	if opts == nil {
		opts = &models.EmployerTransactionListOptions{}
	}

	for {
		opts.Page = page
		opts.Limit = limit

		response, err := s.GetEmployerTransactions(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get transactions page %d: %w", page, err)
		}

		allTransactions = append(allTransactions, response.Results...)

		// Check if we have more pages
		if len(response.Results) < limit {
			break
		}
		page++
	}

	return allTransactions, nil
}

// GetTransactionsByEmployee retrieves all transactions for a specific employee
func (s *TransactionService) GetTransactionsByEmployee(ctx context.Context, employeeID string) ([]models.Transaction, error) {
	opts := &models.TransactionListOptions{
		EmployeeID: employeeID,
		Limit:      1000, // Get all transactions for employee
	}

	result, err := s.GetEmployeeTransactionHistory(ctx, employeeID, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for employee %s: %w", employeeID, err)
	}

	return result.Transactions, nil
}

// GetPendingTransactions retrieves all pending transactions
func (s *TransactionService) GetPendingTransactions(ctx context.Context) ([]models.EmployerTransaction, error) {
	opts := &models.EmployerTransactionListOptions{
		Status: "pending",
	}

	result, err := s.GetAllEmployerTransactions(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending transactions: %w", err)
	}

	return result, nil
}

// GetTransactionsByDateRange retrieves transactions within a date range
func (s *TransactionService) GetTransactionsByDateRange(ctx context.Context, startDate, endDate string) ([]models.EmployerTransaction, error) {
	opts := &models.EmployerTransactionListOptions{
		StartDate: startDate,
		EndDate:   endDate,
	}

	result, err := s.GetAllEmployerTransactions(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for date range %s to %s: %w", startDate, endDate, err)
	}

	return result, nil
}

// CreateAdvanceTransaction creates an advance transaction for an employee
func (s *TransactionService) CreateAdvanceTransaction(ctx context.Context, employeeID string, amount float64, description string) (*models.Transaction, error) {
	req := models.TransactionRequest{
		EmployeeID:  employeeID,
		Amount:      amount,
		Type:        "advance",
		Description: description,
	}

	return s.CreateEmployeeTransaction(ctx, req)
}

// CreateRepaymentTransaction creates a repayment transaction for an employee
func (s *TransactionService) CreateRepaymentTransaction(ctx context.Context, employeeID string, amount float64, description string) (*models.Transaction, error) {
	req := models.TransactionRequest{
		EmployeeID:  employeeID,
		Amount:      amount,
		Type:        "repayment",
		Description: description,
	}

	return s.CreateEmployeeTransaction(ctx, req)
}