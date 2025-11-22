package models

// APIResponse represents the standard API response structure
type APIResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// PaginatedData represents paginated response data
type PaginatedData struct {
	Total   int         `json:"total"`
	Results interface{} `json:"results"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the login response data
type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
	User      User   `json:"user"`
}

// User represents user information
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
	Details    string `json:"details,omitempty"`
}