package models

// EmployeeLoginRequest represents an employee login request
type EmployeeLoginRequest struct {
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	EmiratesID string `json:"emiratesId" validate:"required"`
}

// EmployerLoginRequest represents an employer login request  
type EmployerLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ThirdPartyLoginRequest represents a third-party system login request
type ThirdPartyLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	ClientID string `json:"clientId,omitempty"`
	Scope    string `json:"scope,omitempty"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token        string `json:"token"`
	TokenType    string `json:"tokenType,omitempty"`
	ExpiresIn    int    `json:"expiresIn,omitempty"`
	ExpiresAt    string `json:"expiresAt,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	User         AuthUser `json:"user,omitempty"`
}

// AuthUser represents authenticated user information
type AuthUser struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	FirstName        string `json:"firstName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	Role             string `json:"role"`
	OrganizationID   string `json:"organizationId,omitempty"`
	OrganizationName string `json:"organizationName,omitempty"`
	Permissions      []string `json:"permissions,omitempty"`
	IsActive         bool   `json:"isActive"`
	LastLoginAt      string `json:"lastLoginAt,omitempty"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// LogoutRequest represents a logout request
type LogoutRequest struct {
	Token string `json:"token,omitempty"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
}

// ResetPasswordRequest represents a password reset request
type ResetPasswordRequest struct {
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	EmiratesID string `json:"emiratesId,omitempty"`
}

// ResetPasswordConfirmRequest represents password reset confirmation
type ResetPasswordConfirmRequest struct {
	Token           string `json:"token" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
}

// MFASetupRequest represents MFA setup request
type MFASetupRequest struct {
	Method string `json:"method" validate:"required,oneof=sms email totp"`
	Phone  string `json:"phone,omitempty"`
	Email  string `json:"email,omitempty"`
}

// MFAVerificationRequest represents MFA verification request
type MFAVerificationRequest struct {
	Token string `json:"token" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

// MFAResponse represents MFA setup response
type MFAResponse struct {
	Secret    string `json:"secret,omitempty"`
	QRCode    string `json:"qrCode,omitempty"`
	BackupCodes []string `json:"backupCodes,omitempty"`
	Method    string `json:"method"`
	IsEnabled bool   `json:"isEnabled"`
}

// SessionInfo represents current session information
type SessionInfo struct {
	User         AuthUser `json:"user"`
	Token        string   `json:"token"`
	ExpiresAt    string   `json:"expiresAt"`
	LoginTime    string   `json:"loginTime"`
	LastActivity string   `json:"lastActivity"`
	IPAddress    string   `json:"ipAddress,omitempty"`
	UserAgent    string   `json:"userAgent,omitempty"`
}