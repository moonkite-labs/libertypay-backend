package models

import "time"

// Employee represents an employee in the system
type Employee struct {
	ID              string    `json:"id,omitempty"`
	EmployeeCode    string    `json:"employeeCode" validate:"required"`
	FirstName       string    `json:"firstName" validate:"required"`
	LastName        string    `json:"lastName" validate:"required"`
	Department      string    `json:"department" validate:"required"`
	Designation     string    `json:"designation" validate:"required"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email" validate:"required,email"`
	DOB             string    `json:"dob" validate:"required"` // Format: YYYY-MM-DD
	DateOfJoining   string    `json:"dateOfJoining" validate:"required"` // Format: YYYY-MM-DD
	AccountTitle    string    `json:"accountTitle" validate:"required"`
	AccountNumber   string    `json:"accountNumber" validate:"required"`
	NetSalary       string    `json:"netSalary" validate:"required"`
	EmiratesID      string    `json:"emiratesId" validate:"required"`
	Gender          string    `json:"gender" validate:"required,oneof=Male Female male female"`
	BankID          string    `json:"bankId" validate:"required,uuid4"`
	PayrollStartDay int       `json:"payrollStartDay" validate:"required,min=1,max=31"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
	UpdatedAt       time.Time `json:"updatedAt,omitempty"`
}

// EmployeesRequest represents a request to add/update multiple employees
type EmployeesRequest struct {
	Employees []Employee `json:"employees" validate:"required,min=1,dive"`
}

// EmployeeListOptions represents query options for listing employees
type EmployeeListOptions struct {
	Page       int    `json:"page,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Search     string `json:"search,omitempty"`
	Department string `json:"department,omitempty"`
	Status     string `json:"status,omitempty"`
}

// EmployeeListResponse represents the response for employee list
type EmployeeListResponse struct {
	Total   int        `json:"total"`
	Results []Employee `json:"results"`
}

// EmployeeResponse represents a single employee response
type EmployeeResponse struct {
	Employee Employee `json:"employee"`
}

