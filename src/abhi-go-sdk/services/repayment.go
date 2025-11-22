package services

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"abhi-go-sdk/client"
	"abhi-go-sdk/models"
)

// RepaymentService handles repayment-related API operations
type RepaymentService struct {
	client *client.Client
}

// NewRepaymentService creates a new repayment service
func NewRepaymentService(client *client.Client) *RepaymentService {
	return &RepaymentService{
		client: client,
	}
}

// Create creates a new repayment
func (s *RepaymentService) Create(ctx context.Context, req models.CreateRepaymentRequest) (*models.RepaymentResponse, error) {
	var result models.RepaymentResponse
	err := s.client.POST(ctx, "/repayments", req, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create repayment: %w", err)
	}

	return &result, nil
}

// GetOutstandingBalance retrieves outstanding balance information
func (s *RepaymentService) GetOutstandingBalance(ctx context.Context, opts *models.OutstandingBalanceListOptions) (*models.OutstandingBalanceListResponse, error) {
	query := url.Values{}
	
	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.EmployeeID != "" {
			query.Set("employeeId", opts.EmployeeID)
		}
		if opts.EmployeeCode != "" {
			query.Set("employeeCode", opts.EmployeeCode)
		}
		if opts.Department != "" {
			query.Set("department", opts.Department)
		}
		if opts.MinAmount > 0 {
			query.Set("minAmount", strconv.FormatFloat(opts.MinAmount, 'f', 2, 64))
		}
		if opts.MaxAmount > 0 {
			query.Set("maxAmount", strconv.FormatFloat(opts.MaxAmount, 'f', 2, 64))
		}
		if opts.Overdue {
			query.Set("overdue", "true")
		}
	}

	var result models.OutstandingBalanceListResponse
	err := s.client.GETWithQuery(ctx, "/repayments/outstanding", query, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get outstanding balance: %w", err)
	}

	return &result, nil
}

// GetEmployeeOutstandingBalance retrieves outstanding balance for a specific employee
func (s *RepaymentService) GetEmployeeOutstandingBalance(ctx context.Context, employeeID string) (*models.OutstandingBalance, error) {
	opts := &models.OutstandingBalanceListOptions{
		EmployeeID: employeeID,
		Limit:      1,
	}

	result, err := s.GetOutstandingBalance(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee outstanding balance: %w", err)
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("no outstanding balance found for employee %s", employeeID)
	}

	return &result.Results[0], nil
}

// ListRepayments retrieves a paginated list of repayments
func (s *RepaymentService) ListRepayments(ctx context.Context, opts *models.RepaymentListOptions) (*models.RepaymentListResponse, error) {
	query := url.Values{}
	
	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.EmployeeID != "" {
			query.Set("employeeId", opts.EmployeeID)
		}
		if opts.Status != "" {
			query.Set("status", opts.Status)
		}
		if opts.StartDate != "" {
			query.Set("startDate", opts.StartDate)
		}
		if opts.EndDate != "" {
			query.Set("endDate", opts.EndDate)
		}
		if opts.ClientRepaymentReferenceNumber != "" {
			query.Set("clientRepaymentReferenceNumber", opts.ClientRepaymentReferenceNumber)
		}
		if opts.MinAmount > 0 {
			query.Set("minAmount", strconv.FormatFloat(opts.MinAmount, 'f', 2, 64))
		}
		if opts.MaxAmount > 0 {
			query.Set("maxAmount", strconv.FormatFloat(opts.MaxAmount, 'f', 2, 64))
		}
	}

	var result models.RepaymentListResponse
	err := s.client.GETWithQuery(ctx, "/repayments", query, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to list repayments: %w", err)
	}

	return &result, nil
}

// GetRepaymentByID retrieves a specific repayment by ID
func (s *RepaymentService) GetRepaymentByID(ctx context.Context, repaymentID string) (*models.Repayment, error) {
	var result models.Repayment
	endpoint := fmt.Sprintf("/repayments/%s", repaymentID)
	
	err := s.client.GET(ctx, endpoint, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get repayment %s: %w", repaymentID, err)
	}

	return &result, nil
}

// GetRepaymentByReference retrieves a repayment by client reference number
func (s *RepaymentService) GetRepaymentByReference(ctx context.Context, referenceNumber string) (*models.Repayment, error) {
	opts := &models.RepaymentListOptions{
		ClientRepaymentReferenceNumber: referenceNumber,
		Limit: 1,
	}

	result, err := s.ListRepayments(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get repayment by reference %s: %w", referenceNumber, err)
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("repayment with reference %s not found", referenceNumber)
	}

	return &result.Results[0], nil
}

// GetEmployeeRepayments retrieves all repayments for a specific employee
func (s *RepaymentService) GetEmployeeRepayments(ctx context.Context, employeeID string) ([]models.Repayment, error) {
	var allRepayments []models.Repayment
	page := 1
	limit := 100

	for {
		opts := &models.RepaymentListOptions{
			EmployeeID: employeeID,
			Page:       page,
			Limit:      limit,
		}

		response, err := s.ListRepayments(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get employee repayments page %d: %w", page, err)
		}

		allRepayments = append(allRepayments, response.Results...)

		// Check if we have more pages
		if len(response.Results) < limit {
			break
		}
		page++
	}

	return allRepayments, nil
}

// GetOverdueBalances retrieves all overdue outstanding balances
func (s *RepaymentService) GetOverdueBalances(ctx context.Context) ([]models.OutstandingBalance, error) {
	opts := &models.OutstandingBalanceListOptions{
		Overdue: true,
		Limit:   1000, // Get all overdue balances
	}

	result, err := s.GetOutstandingBalance(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue balances: %w", err)
	}

	return result.Results, nil
}

// GetRepaymentsByDateRange retrieves repayments within a date range
func (s *RepaymentService) GetRepaymentsByDateRange(ctx context.Context, startDate, endDate string) ([]models.Repayment, error) {
	var allRepayments []models.Repayment
	page := 1
	limit := 100

	for {
		opts := &models.RepaymentListOptions{
			StartDate: startDate,
			EndDate:   endDate,
			Page:      page,
			Limit:     limit,
		}

		response, err := s.ListRepayments(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get repayments for date range %s to %s page %d: %w", startDate, endDate, page, err)
		}

		allRepayments = append(allRepayments, response.Results...)

		// Check if we have more pages
		if len(response.Results) < limit {
			break
		}
		page++
	}

	return allRepayments, nil
}

// GetRepaymentsByStatus retrieves repayments by status
func (s *RepaymentService) GetRepaymentsByStatus(ctx context.Context, status string) ([]models.Repayment, error) {
	var allRepayments []models.Repayment
	page := 1
	limit := 100

	for {
		opts := &models.RepaymentListOptions{
			Status: status,
			Page:   page,
			Limit:  limit,
		}

		response, err := s.ListRepayments(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get repayments with status %s page %d: %w", status, page, err)
		}

		allRepayments = append(allRepayments, response.Results...)

		// Check if we have more pages
		if len(response.Results) < limit {
			break
		}
		page++
	}

	return allRepayments, nil
}

// CreateEmployeeRepayment creates a repayment for a specific employee
func (s *RepaymentService) CreateEmployeeRepayment(ctx context.Context, employeeID string, amount float64, referenceNumber, description string) (*models.RepaymentResponse, error) {
	req := models.CreateRepaymentRequest{
		Amount:                         amount,
		ClientRepaymentReferenceNumber: referenceNumber,
		EmployeeID:                     employeeID,
		Description:                    description,
	}

	return s.Create(ctx, req)
}

// CreateTransactionRepayment creates a repayment for a specific transaction
func (s *RepaymentService) CreateTransactionRepayment(ctx context.Context, transactionID string, amount float64, referenceNumber, description string) (*models.RepaymentResponse, error) {
	req := models.CreateRepaymentRequest{
		Amount:                         amount,
		ClientRepaymentReferenceNumber: referenceNumber,
		TransactionID:                  transactionID,
		Description:                    description,
	}

	return s.Create(ctx, req)
}

// ValidateRepayment validates repayment data before creation
func (s *RepaymentService) ValidateRepayment(req models.CreateRepaymentRequest) error {
	if req.Amount <= 0 {
		return fmt.Errorf("repayment amount must be greater than 0")
	}
	if req.ClientRepaymentReferenceNumber == "" {
		return fmt.Errorf("client repayment reference number is required")
	}
	if req.EmployeeID == "" && req.TransactionID == "" {
		return fmt.Errorf("either employee ID or transaction ID must be provided")
	}

	return nil
}

// GetOutstandingBalanceSummary returns summary statistics for outstanding balances
func (s *RepaymentService) GetOutstandingBalanceSummary(ctx context.Context) (*models.OutstandingBalanceSummary, error) {
	result, err := s.GetOutstandingBalance(ctx, &models.OutstandingBalanceListOptions{
		Limit: 1000, // Get all outstanding balances
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get outstanding balance summary: %w", err)
	}

	return &result.Summary, nil
}