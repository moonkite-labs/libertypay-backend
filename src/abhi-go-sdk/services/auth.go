package services

import (
	"context"
	"fmt"

	"abhi-go-sdk/client"
	"abhi-go-sdk/models"
)

// AuthService handles authentication-related API operations
type AuthService struct {
	client *client.Client
}

// NewAuthService creates a new authentication service
func NewAuthService(client *client.Client) *AuthService {
	return &AuthService{
		client: client,
	}
}

// EmployeeLogin authenticates an employee with username, password, and Emirates ID
func (s *AuthService) EmployeeLogin(ctx context.Context, req models.EmployeeLoginRequest) (*models.AuthResponse, error) {
	var result models.AuthResponse
	err := s.client.POST(ctx, "/auth/employee-login", req, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to login employee: %w", err)
	}

	return &result, nil
}

// EmployerLogin authenticates an employer with username and password
func (s *AuthService) EmployerLogin(ctx context.Context, req models.EmployerLoginRequest) (*models.AuthResponse, error) {
	var result models.AuthResponse
	err := s.client.POST(ctx, "/auth/employer-login", req, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to login employer: %w", err)
	}

	return &result, nil
}

// ThirdPartyLogin authenticates a third-party system
func (s *AuthService) ThirdPartyLogin(ctx context.Context, req models.ThirdPartyLoginRequest) (*models.AuthResponse, error) {
	var result models.AuthResponse
	err := s.client.POST(ctx, "/auth/login", req, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to login third-party: %w", err)
	}

	return &result, nil
}

// RefreshToken refreshes an existing authentication token
func (s *AuthService) RefreshToken(ctx context.Context, req models.RefreshTokenRequest) (*models.AuthResponse, error) {
	var result models.AuthResponse
	err := s.client.POST(ctx, "/auth/refresh", req, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &result, nil
}

// Logout invalidates the current session
func (s *AuthService) Logout(ctx context.Context, req models.LogoutRequest) error {
	err := s.client.POST(ctx, "/auth/logout", req, nil)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}

// GetCurrentUser retrieves information about the currently authenticated user
func (s *AuthService) GetCurrentUser(ctx context.Context) (*models.AuthUser, error) {
	var result models.AuthUser
	err := s.client.GET(ctx, "/auth/me", &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	return &result, nil
}

// GetSessionInfo retrieves current session information
func (s *AuthService) GetSessionInfo(ctx context.Context) (*models.SessionInfo, error) {
	var result models.SessionInfo
	err := s.client.GET(ctx, "/auth/session", &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get session info: %w", err)
	}

	return &result, nil
}

// ChangePassword changes the password for the current user
func (s *AuthService) ChangePassword(ctx context.Context, req models.ChangePasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return fmt.Errorf("new password and confirm password do not match")
	}

	err := s.client.POST(ctx, "/auth/change-password", req, nil)
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}

// RequestPasswordReset initiates a password reset process
func (s *AuthService) RequestPasswordReset(ctx context.Context, req models.ResetPasswordRequest) error {
	err := s.client.POST(ctx, "/auth/reset-password", req, nil)
	if err != nil {
		return fmt.Errorf("failed to request password reset: %w", err)
	}

	return nil
}

// ConfirmPasswordReset confirms and completes the password reset process
func (s *AuthService) ConfirmPasswordReset(ctx context.Context, req models.ResetPasswordConfirmRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return fmt.Errorf("new password and confirm password do not match")
	}

	err := s.client.POST(ctx, "/auth/reset-password/confirm", req, nil)
	if err != nil {
		return fmt.Errorf("failed to confirm password reset: %w", err)
	}

	return nil
}

// SetupMFA sets up multi-factor authentication for the current user
func (s *AuthService) SetupMFA(ctx context.Context, req models.MFASetupRequest) (*models.MFAResponse, error) {
	var result models.MFAResponse
	err := s.client.POST(ctx, "/auth/mfa/setup", req, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to setup MFA: %w", err)
	}

	return &result, nil
}

// VerifyMFA verifies MFA during login or setup
func (s *AuthService) VerifyMFA(ctx context.Context, req models.MFAVerificationRequest) (*models.AuthResponse, error) {
	var result models.AuthResponse
	err := s.client.POST(ctx, "/auth/mfa/verify", req, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to verify MFA: %w", err)
	}

	return &result, nil
}

// DisableMFA disables multi-factor authentication for the current user
func (s *AuthService) DisableMFA(ctx context.Context) error {
	err := s.client.POST(ctx, "/auth/mfa/disable", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to disable MFA: %w", err)
	}

	return nil
}

// GetMFAStatus retrieves the current MFA status
func (s *AuthService) GetMFAStatus(ctx context.Context) (*models.MFAResponse, error) {
	var result models.MFAResponse
	err := s.client.GET(ctx, "/auth/mfa/status", &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get MFA status: %w", err)
	}

	return &result, nil
}

// ValidateToken validates the current authentication token
func (s *AuthService) ValidateToken(ctx context.Context) (bool, error) {
	err := s.client.GET(ctx, "/auth/validate", nil)
	if err != nil {
		// If validation fails, token is invalid
		return false, nil
	}

	return true, nil
}

// Convenience Methods

// LoginEmployee is a convenience method for employee login
func (s *AuthService) LoginEmployee(ctx context.Context, username, password, emiratesID string) (*models.AuthResponse, error) {
	req := models.EmployeeLoginRequest{
		Username:   username,
		Password:   password,
		EmiratesID: emiratesID,
	}

	return s.EmployeeLogin(ctx, req)
}

// LoginEmployer is a convenience method for employer login
func (s *AuthService) LoginEmployer(ctx context.Context, username, password string) (*models.AuthResponse, error) {
	req := models.EmployerLoginRequest{
		Username: username,
		Password: password,
	}

	return s.EmployerLogin(ctx, req)
}

// LoginThirdParty is a convenience method for third-party login
func (s *AuthService) LoginThirdParty(ctx context.Context, username, password string) (*models.AuthResponse, error) {
	req := models.ThirdPartyLoginRequest{
		Username: username,
		Password: password,
	}

	return s.ThirdPartyLogin(ctx, req)
}

// LogoutCurrentSession is a convenience method to logout the current session
func (s *AuthService) LogoutCurrentSession(ctx context.Context) error {
	return s.Logout(ctx, models.LogoutRequest{})
}

// ValidateCredentials validates login credentials without creating a session
func (s *AuthService) ValidateCredentials(ctx context.Context, loginType string, credentials interface{}) (bool, error) {
	switch loginType {
	case "employee":
		req, ok := credentials.(models.EmployeeLoginRequest)
		if !ok {
			return false, fmt.Errorf("invalid employee credentials format")
		}
		_, err := s.EmployeeLogin(ctx, req)
		return err == nil, err
	case "employer":
		req, ok := credentials.(models.EmployerLoginRequest)
		if !ok {
			return false, fmt.Errorf("invalid employer credentials format")
		}
		_, err := s.EmployerLogin(ctx, req)
		return err == nil, err
	case "third-party":
		req, ok := credentials.(models.ThirdPartyLoginRequest)
		if !ok {
			return false, fmt.Errorf("invalid third-party credentials format")
		}
		_, err := s.ThirdPartyLogin(ctx, req)
		return err == nil, err
	default:
		return false, fmt.Errorf("unsupported login type: %s", loginType)
	}
}