package services

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"abhi-go-sdk/client"
	"abhi-go-sdk/models"
)

// OrganizationService handles organization-related API operations
type OrganizationService struct {
	client *client.Client
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(client *client.Client) *OrganizationService {
	return &OrganizationService{
		client: client,
	}
}

// List retrieves a paginated list of sub-organizations
func (s *OrganizationService) List(ctx context.Context, opts *models.OrganizationListOptions) (*models.OrganizationListResponse, error) {
	query := url.Values{}
	
	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.From != "" {
			query.Set("from", opts.From)
		}
		if opts.To != "" {
			query.Set("to", opts.To)
		}
		if opts.ShowInactive {
			query.Set("showInactive", "true")
		}
		if opts.Column != "" {
			query.Set("column", opts.Column)
		}
		if opts.Order != "" {
			query.Set("order", opts.Order)
		}
	}

	var result models.OrganizationListResponse
	err := s.client.GETWithQuery(ctx, "/organizations", query, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	return &result, nil
}

// GetAll retrieves all organizations with pagination handling
func (s *OrganizationService) GetAll(ctx context.Context) ([]models.Organization, error) {
	var allOrganizations []models.Organization
	page := 1
	limit := 100

	for {
		opts := &models.OrganizationListOptions{
			Page:  page,
			Limit: limit,
		}

		response, err := s.List(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get organizations page %d: %w", page, err)
		}

		allOrganizations = append(allOrganizations, response.Results...)

		// Check if we have more pages
		if len(response.Results) < limit {
			break
		}
		page++
	}

	return allOrganizations, nil
}

// GetByID retrieves a single organization by ID
func (s *OrganizationService) GetByID(ctx context.Context, organizationID string) (*models.Organization, error) {
	var result models.Organization
	endpoint := fmt.Sprintf("/organizations/%s", organizationID)
	
	err := s.client.GET(ctx, endpoint, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization %s: %w", organizationID, err)
	}

	return &result, nil
}

// Create creates a new sub-organization
func (s *OrganizationService) Create(ctx context.Context, req models.CreateOrganizationRequest) (*models.CreateOrganizationResponse, error) {
	var result models.CreateOrganizationResponse
	err := s.client.POST(ctx, "/organizations", req, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	return &result, nil
}

// GetActive retrieves only active organizations
func (s *OrganizationService) GetActive(ctx context.Context, opts *models.OrganizationListOptions) (*models.OrganizationListResponse, error) {
	if opts == nil {
		opts = &models.OrganizationListOptions{}
	}
	opts.ShowInactive = false

	result, err := s.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get active organizations: %w", err)
	}

	return result, nil
}

// GetInactive retrieves only inactive organizations
func (s *OrganizationService) GetInactive(ctx context.Context, opts *models.OrganizationListOptions) (*models.OrganizationListResponse, error) {
	if opts == nil {
		opts = &models.OrganizationListOptions{}
	}
	opts.ShowInactive = true

	result, err := s.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get inactive organizations: %w", err)
	}

	// Filter to only inactive organizations
	var inactiveOrgs []models.Organization
	for _, org := range result.Results {
		if !org.Active {
			inactiveOrgs = append(inactiveOrgs, org)
		}
	}

	return &models.OrganizationListResponse{
		Total:   len(inactiveOrgs),
		Results: inactiveOrgs,
	}, nil
}

// GetByIndustry retrieves organizations by industry
func (s *OrganizationService) GetByIndustry(ctx context.Context, industry string) ([]models.Organization, error) {
	allOrgs, err := s.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations by industry %s: %w", industry, err)
	}

	var industryOrgs []models.Organization
	for _, org := range allOrgs {
		if org.Industry == industry {
			industryOrgs = append(industryOrgs, org)
		}
	}

	return industryOrgs, nil
}

// GetByDateRange retrieves organizations created within a date range
func (s *OrganizationService) GetByDateRange(ctx context.Context, startDate, endDate string) ([]models.Organization, error) {
	opts := &models.OrganizationListOptions{
		From:  startDate,
		To:    endDate,
		Limit: 1000, // Get all organizations in range
	}

	result, err := s.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations for date range %s to %s: %w", startDate, endDate, err)
	}

	return result.Results, nil
}

// Search searches for organizations by name
func (s *OrganizationService) Search(ctx context.Context, searchTerm string, limit int) ([]models.Organization, error) {
	if limit <= 0 {
		limit = 50
	}

	// Get all organizations and filter by name
	allOrgs, err := s.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search organizations: %w", err)
	}

	var matchedOrgs []models.Organization
	count := 0
	for _, org := range allOrgs {
		if count >= limit {
			break
		}
		// Simple case-insensitive name search
		if len(org.Name) >= len(searchTerm) {
			for i := 0; i <= len(org.Name)-len(searchTerm); i++ {
				if org.Name[i:i+len(searchTerm)] == searchTerm {
					matchedOrgs = append(matchedOrgs, org)
					count++
					break
				}
			}
		}
	}

	return matchedOrgs, nil
}

// GetSortedByName retrieves organizations sorted by name
func (s *OrganizationService) GetSortedByName(ctx context.Context, ascending bool) ([]models.Organization, error) {
	order := "DESC"
	if ascending {
		order = "ASC"
	}

	opts := &models.OrganizationListOptions{
		Column: "organizations.name",
		Order:  order,
		Limit:  1000, // Get all organizations
	}

	result, err := s.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations sorted by name: %w", err)
	}

	return result.Results, nil
}

// GetSortedByCreationDate retrieves organizations sorted by creation date
func (s *OrganizationService) GetSortedByCreationDate(ctx context.Context, ascending bool) ([]models.Organization, error) {
	order := "DESC"
	if ascending {
		order = "ASC"
	}

	opts := &models.OrganizationListOptions{
		Column: "organizations.createdAt",
		Order:  order,
		Limit:  1000, // Get all organizations
	}

	result, err := s.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations sorted by creation date: %w", err)
	}

	return result.Results, nil
}

// ValidateOrganization validates organization data before creation
func (s *OrganizationService) ValidateOrganization(req models.CreateOrganizationRequest) error {
	if req.Name == "" {
		return fmt.Errorf("organization name is required")
	}
	if req.Industry == "" {
		return fmt.Errorf("industry is required")
	}
	if req.BusinessTypeID == "" {
		return fmt.Errorf("business type ID is required")
	}
	if req.Address == "" {
		return fmt.Errorf("address is required")
	}
	if req.City == "" {
		return fmt.Errorf("city is required")
	}
	if req.ManagementAlias == "" {
		return fmt.Errorf("management alias is required")
	}
	if len(req.ManagementAlias) < 4 || len(req.ManagementAlias) > 100 {
		return fmt.Errorf("management alias must be between 4 and 100 characters")
	}
	if req.CreditLimit <= 0 {
		return fmt.Errorf("credit limit must be greater than 0")
	}
	if req.PayrollStartDay < 0 || req.PayrollStartDay > 31 {
		return fmt.Errorf("payroll start day must be between 1 and 31")
	}

	return nil
}

// GetStatistics returns organization statistics
func (s *OrganizationService) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
	allOrgs, err := s.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization statistics: %w", err)
	}

	stats := map[string]interface{}{
		"total":      len(allOrgs),
		"active":     0,
		"inactive":   0,
		"industries": make(map[string]int),
	}

	industries := stats["industries"].(map[string]int)

	for _, org := range allOrgs {
		if org.Active {
			stats["active"] = stats["active"].(int) + 1
		} else {
			stats["inactive"] = stats["inactive"].(int) + 1
		}

		industries[org.Industry]++
	}

	return stats, nil
}