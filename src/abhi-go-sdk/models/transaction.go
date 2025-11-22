package models

import (
	"time"
	"abhi-go-sdk/errors"
)

// Transaction represents a transaction in the system
type Transaction struct {
	ID                string    `json:"id,omitempty"`
	EmployeeID        string    `json:"employeeId" validate:"required"`
	Amount            float64   `json:"amount" validate:"required,gt=0"`
	Type              string    `json:"type" validate:"required,oneof=advance repayment"`
	Status            string    `json:"status,omitempty"`
	Description       string    `json:"description,omitempty"`
	RequestedAt       time.Time `json:"requestedAt,omitempty"`
	ProcessedAt       time.Time `json:"processedAt,omitempty"`
	DueDate           string    `json:"dueDate,omitempty"`
	RepaymentAmount   float64   `json:"repaymentAmount,omitempty"`
	InterestRate      float64   `json:"interestRate,omitempty"`
	ProcessingFee     float64   `json:"processingFee,omitempty"`
	TransactionRef    string    `json:"transactionRef,omitempty"`
	BankTransactionID string    `json:"bankTransactionId,omitempty"`
	Reason            string    `json:"reason,omitempty"`
}

// TransactionRequest represents a transaction request
type TransactionRequest struct {
	EmployeeID  string  `json:"employeeId" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Type        string  `json:"type" validate:"required,oneof=advance repayment"`
	Description string  `json:"description,omitempty"`
	DueDate     string  `json:"dueDate,omitempty"`
}

// TransactionListOptions represents query options for listing transactions
type TransactionListOptions struct {
	Page       int    `json:"page,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	EmployeeID string `json:"employeeId,omitempty"`
	Status     string `json:"status,omitempty"`
	Type       string `json:"type,omitempty"`
	StartDate  string `json:"startDate,omitempty"`
	EndDate    string `json:"endDate,omitempty"`
}

// TransactionListResponse represents the response for transaction list
type TransactionListResponse struct {
	Total   int           `json:"total"`
	Results []Transaction `json:"results"`
}

// TransactionHistoryResponse represents employee transaction history
type TransactionHistoryResponse struct {
	EmployeeID   string        `json:"employeeId"`
	TotalCount   int           `json:"totalCount"`
	Transactions []Transaction `json:"transactions"`
}

// MonthlyBalance represents monthly balance information
type MonthlyBalance struct {
	Month           string  `json:"month"`
	Year            int     `json:"year"`
	GrossSalary     float64 `json:"grossSalary"`
	NetSalary       float64 `json:"netSalary"`
	AvailableAmount float64 `json:"availableAmount"`
	UsedAmount      float64 `json:"usedAmount"`
	PendingAmount   float64 `json:"pendingAmount"`
	LastUpdated     string  `json:"lastUpdated"`
}

// MonthlyBalanceResponse represents monthly balance response
type MonthlyBalanceResponse struct {
	EmployeeID string         `json:"employeeId"`
	Balance    MonthlyBalance `json:"balance"`
}

// TransactionValidationRequest represents transaction validation request
type TransactionValidationRequest struct {
	EmployeeID string  `json:"employeeId" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}

// TransactionValidationResponse represents transaction validation response
type TransactionValidationResponse struct {
	IsValid          bool    `json:"isValid"`
	MaxAmount        float64 `json:"maxAmount"`
	AvailableAmount  float64 `json:"availableAmount"`
	Message          string  `json:"message"`
	ValidationErrors []errors.ValidationError `json:"validationErrors,omitempty"`
}

// TransactionStatusResponse represents transaction status response
type TransactionStatusResponse struct {
	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
	Message       string `json:"message"`
	LastUpdated   string `json:"lastUpdated"`
}

// EmployerTransactionListOptions represents query options for employer transaction listing
type EmployerTransactionListOptions struct {
	Page         int    `json:"page,omitempty"`
	Limit        int    `json:"limit,omitempty"`
	Status       string `json:"status,omitempty"`
	Type         string `json:"type,omitempty"`
	StartDate    string `json:"startDate,omitempty"`
	EndDate      string `json:"endDate,omitempty"`
	EmployeeCode string `json:"employeeCode,omitempty"`
	Department   string `json:"department,omitempty"`
}

// EmployerTransactionResponse represents employer view of transactions
type EmployerTransactionResponse struct {
	Total   int                     `json:"total"`
	Results []EmployerTransaction   `json:"results"`
}

// EmployerTransaction represents transaction from employer perspective
type EmployerTransaction struct {
	ID              string  `json:"id"`
	EmployeeID      string  `json:"employeeId"`
	EmployeeCode    string  `json:"employeeCode"`
	EmployeeName    string  `json:"employeeName"`
	Department      string  `json:"department"`
	Amount          float64 `json:"amount"`
	Type            string  `json:"type"`
	Status          string  `json:"status"`
	RequestedAt     string  `json:"requestedAt"`
	ProcessedAt     string  `json:"processedAt"`
	DueDate         string  `json:"dueDate"`
	RepaymentAmount float64 `json:"repaymentAmount"`
}

// ValidationQuestion represents a validation question
type ValidationQuestion struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Type     string `json:"type"` // text, multiple_choice, yes_no
	Required bool   `json:"required"`
	Options  []string `json:"options,omitempty"`
}

// ValidationAnswer represents an answer to a validation question
type ValidationAnswer struct {
	QuestionID string `json:"questionId" validate:"required"`
	Answer     string `json:"answer" validate:"required"`
}

// ValidationQuestionsRequest represents validation questions request
type ValidationQuestionsRequest struct {
	TransactionID string `json:"transactionId" validate:"required"`
}

// ValidationQuestionsResponse represents validation questions response
type ValidationQuestionsResponse struct {
	TransactionID string               `json:"transactionId"`
	Questions     []ValidationQuestion `json:"questions"`
}

// ValidationAnswersRequest represents validation answers submission
type ValidationAnswersRequest struct {
	TransactionID string             `json:"transactionId" validate:"required"`
	Answers       []ValidationAnswer `json:"answers" validate:"required,dive"`
}

// ValidationAnswersResponse represents validation answers response
type ValidationAnswersResponse struct {
	TransactionID string `json:"transactionId"`
	IsValid       bool   `json:"isValid"`
	Message       string `json:"message"`
}