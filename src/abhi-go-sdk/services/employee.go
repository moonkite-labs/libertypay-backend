package services

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"abhi-go-sdk/client"
	"abhi-go-sdk/models"
)

// EmployeeService handles employee-related API operations
type EmployeeService struct {
	client *client.Client
}

// NewEmployeeService creates a new employee service
func NewEmployeeService(client *client.Client) *EmployeeService {
	return &EmployeeService{
		client: client,
	}
}

// List retrieves a paginated list of employees
func (s *EmployeeService) List(ctx context.Context, opts *models.EmployeeListOptions) (*models.EmployeeListResponse, error) {
	query := url.Values{}
	
	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Search != "" {
			query.Set("search", opts.Search)
		}
		if opts.Department != "" {
			query.Set("department", opts.Department)
		}
		if opts.Status != "" {
			query.Set("status", opts.Status)
		}
	}

	var result models.EmployeeListResponse
	err := s.client.GETWithQuery(ctx, "/employees", query, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to list employees: %w", err)
	}

	return &result, nil
}

// GetAll retrieves all employees with pagination handling
func (s *EmployeeService) GetAll(ctx context.Context) ([]models.Employee, error) {
	var allEmployees []models.Employee
	page := 1
	limit := 100

	for {
		opts := &models.EmployeeListOptions{
			Page:  page,
			Limit: limit,
		}

		response, err := s.List(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get employees page %d: %w", page, err)
		}

		allEmployees = append(allEmployees, response.Results...)

		// Check if we have more pages
		if len(response.Results) < limit {
			break
		}
		page++
	}

	return allEmployees, nil
}

// GetByID retrieves a single employee by ID
func (s *EmployeeService) GetByID(ctx context.Context, employeeID string) (*models.Employee, error) {
	var result models.Employee
	endpoint := fmt.Sprintf("/employees/%s", employeeID)
	
	err := s.client.GET(ctx, endpoint, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee %s: %w", employeeID, err)
	}

	return &result, nil
}

// GetByEmployeeCode retrieves a single employee by employee code
func (s *EmployeeService) GetByEmployeeCode(ctx context.Context, employeeCode string) (*models.Employee, error) {
	opts := &models.EmployeeListOptions{
		Search: employeeCode,
		Limit:  1,
	}

	result, err := s.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search for employee code %s: %w", employeeCode, err)
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("employee with code %s not found", employeeCode)
	}

	// Find exact match
	for _, emp := range result.Results {
		if emp.EmployeeCode == employeeCode {
			return &emp, nil
		}
	}

	return nil, fmt.Errorf("employee with code %s not found", employeeCode)
}

// Create adds new employees to the system
func (s *EmployeeService) Create(ctx context.Context, employees []models.Employee) error {
	request := models.EmployeesRequest{
		Employees: employees,
	}

	err := s.client.POST(ctx, "/employees", request, nil)
	if err != nil {
		return fmt.Errorf("failed to create employees: %w", err)
	}

	return nil
}

// CreateSingle adds a single employee to the system
func (s *EmployeeService) CreateSingle(ctx context.Context, employee models.Employee) error {
	return s.Create(ctx, []models.Employee{employee})
}

// Update updates existing employees
func (s *EmployeeService) Update(ctx context.Context, employees []models.Employee) error {
	request := models.EmployeesRequest{
		Employees: employees,
	}

	err := s.client.PUT(ctx, "/employees", request, nil)
	if err != nil {
		return fmt.Errorf("failed to update employees: %w", err)
	}

	return nil
}

// UpdateSingle updates a single employee
func (s *EmployeeService) UpdateSingle(ctx context.Context, employee models.Employee) error {
	return s.Update(ctx, []models.Employee{employee})
}

// Delete removes an employee from the system
func (s *EmployeeService) Delete(ctx context.Context, employeeID string) error {
	endpoint := fmt.Sprintf("/employees/%s", employeeID)
	
	err := s.client.DELETE(ctx, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to delete employee %s: %w", employeeID, err)
	}

	return nil
}

// Search searches for employees based on criteria
func (s *EmployeeService) Search(ctx context.Context, searchTerm string, limit int) ([]models.Employee, error) {
	if limit <= 0 {
		limit = 50
	}

	opts := &models.EmployeeListOptions{
		Search: searchTerm,
		Limit:  limit,
	}

	result, err := s.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search employees: %w", err)
	}

	return result.Results, nil
}

// GetByDepartment retrieves employees by department
func (s *EmployeeService) GetByDepartment(ctx context.Context, department string) ([]models.Employee, error) {
	opts := &models.EmployeeListOptions{
		Department: department,
		Limit:      1000, // Get all employees in department
	}

	result, err := s.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees by department %s: %w", department, err)
	}

	return result.Results, nil
}

// ValidateEmployee validates employee data before creation/update
func (s *EmployeeService) ValidateEmployee(employee models.Employee) error {
	// This would typically use the validator from the client
	// Additional business logic validation can be added here
	if employee.EmployeeCode == "" {
		return fmt.Errorf("employee code is required")
	}
	if employee.Email == "" {
		return fmt.Errorf("email is required")
	}
	if employee.NetSalary == "" {
		return fmt.Errorf("net salary is required")
	}
	if employee.BankID == "" {
		return fmt.Errorf("bank ID is required")
	}

	return nil
}