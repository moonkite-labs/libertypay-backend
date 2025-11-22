package models

import "time"

// Repayment represents a repayment entity
type Repayment struct {
	ID                             string    `json:"id,omitempty"`
	Amount                         float64   `json:"amount" validate:"required,gt=0"`
	ClientRepaymentReferenceNumber string    `json:"clientRepaymentReferenceNumber" validate:"required"`
	EmployeeID                     string    `json:"employeeId,omitempty"`
	TransactionID                  string    `json:"transactionId,omitempty"`
	Status                         string    `json:"status,omitempty"`
	ProcessedAt                    time.Time `json:"processedAt,omitempty"`
	CreatedAt                      time.Time `json:"createdAt,omitempty"`
	UpdatedAt                      time.Time `json:"updatedAt,omitempty"`
	Description                    string    `json:"description,omitempty"`
	PaymentMethod                  string    `json:"paymentMethod,omitempty"`
	BankTransactionID              string    `json:"bankTransactionId,omitempty"`
}

// CreateRepaymentRequest represents a request to create a repayment
type CreateRepaymentRequest struct {
	Amount                         float64 `json:"amount" validate:"required,gt=0"`
	ClientRepaymentReferenceNumber string  `json:"clientRepaymentReferenceNumber" validate:"required"`
	EmployeeID                     string  `json:"employeeId,omitempty"`
	TransactionID                  string  `json:"transactionId,omitempty"`
	Description                    string  `json:"description,omitempty"`
	PaymentMethod                  string  `json:"paymentMethod,omitempty"`
}

// RepaymentResponse represents the response when creating a repayment
type RepaymentResponse struct {
	Repayment Repayment `json:"repayment"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
}

// OutstandingBalance represents outstanding balance information
type OutstandingBalance struct {
	EmployeeID           string  `json:"employeeId"`
	EmployeeCode         string  `json:"employeeCode,omitempty"`
	EmployeeName         string  `json:"employeeName,omitempty"`
	TotalOutstanding     float64 `json:"totalOutstanding"`
	PrincipalAmount      float64 `json:"principalAmount"`
	InterestAmount       float64 `json:"interestAmount"`
	PenaltyAmount        float64 `json:"penaltyAmount"`
	ProcessingFee        float64 `json:"processingFee"`
	OverdueAmount        float64 `json:"overdueAmount"`
	DaysPastDue          int     `json:"daysPastDue"`
	NextDueDate          string  `json:"nextDueDate,omitempty"`
	LastPaymentDate      string  `json:"lastPaymentDate,omitempty"`
	LastPaymentAmount    float64 `json:"lastPaymentAmount"`
	TransactionHistory   []OutstandingTransaction `json:"transactionHistory,omitempty"`
}

// OutstandingTransaction represents a transaction in the outstanding balance
type OutstandingTransaction struct {
	ID              string  `json:"id"`
	Amount          float64 `json:"amount"`
	Type            string  `json:"type"` // advance, repayment, interest, penalty, fee
	Status          string  `json:"status"`
	Date            string  `json:"date"`
	DueDate         string  `json:"dueDate,omitempty"`
	RemainingAmount float64 `json:"remainingAmount"`
	Description     string  `json:"description,omitempty"`
}

// OutstandingBalanceListOptions represents query options for outstanding balance
type OutstandingBalanceListOptions struct {
	Page         int    `json:"page,omitempty"`
	Limit        int    `json:"limit,omitempty"`
	EmployeeID   string `json:"employeeId,omitempty"`
	EmployeeCode string `json:"employeeCode,omitempty"`
	Department   string `json:"department,omitempty"`
	MinAmount    float64 `json:"minAmount,omitempty"`
	MaxAmount    float64 `json:"maxAmount,omitempty"`
	Overdue      bool   `json:"overdue,omitempty"`
}

// OutstandingBalanceListResponse represents the response for outstanding balance list
type OutstandingBalanceListResponse struct {
	Total   int                  `json:"total"`
	Results []OutstandingBalance `json:"results"`
	Summary OutstandingBalanceSummary `json:"summary,omitempty"`
}

// OutstandingBalanceSummary represents summary statistics for outstanding balances
type OutstandingBalanceSummary struct {
	TotalEmployees       int     `json:"totalEmployees"`
	TotalOutstanding     float64 `json:"totalOutstanding"`
	TotalOverdue         float64 `json:"totalOverdue"`
	AverageOutstanding   float64 `json:"averageOutstanding"`
	EmployeesWithOverdue int     `json:"employeesWithOverdue"`
}

// RepaymentListOptions represents query options for listing repayments
type RepaymentListOptions struct {
	Page                           int    `json:"page,omitempty"`
	Limit                          int    `json:"limit,omitempty"`
	EmployeeID                     string `json:"employeeId,omitempty"`
	Status                         string `json:"status,omitempty"`
	StartDate                      string `json:"startDate,omitempty"`
	EndDate                        string `json:"endDate,omitempty"`
	ClientRepaymentReferenceNumber string `json:"clientRepaymentReferenceNumber,omitempty"`
	MinAmount                      float64 `json:"minAmount,omitempty"`
	MaxAmount                      float64 `json:"maxAmount,omitempty"`
}

// RepaymentListResponse represents the response for repayment list
type RepaymentListResponse struct {
	Total   int         `json:"total"`
	Results []Repayment `json:"results"`
}