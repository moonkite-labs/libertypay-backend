package services

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"abhi-go-sdk/client"
	"abhi-go-sdk/models"
)

// MiscService handles miscellaneous API operations (banks, business types)
type MiscService struct {
	client *client.Client
}

// NewMiscService creates a new miscellaneous service
func NewMiscService(client *client.Client) *MiscService {
	return &MiscService{
		client: client,
	}
}

// Banks Section

// GetBanks retrieves a paginated list of banks
func (s *MiscService) GetBanks(ctx context.Context, opts *models.BankListOptions) (*models.BankListResponse, error) {
	query := url.Values{}
	
	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Country != "" {
			query.Set("country", opts.Country)
		}
		if opts.Active != nil {
			query.Set("active", strconv.FormatBool(*opts.Active))
		}
	}

	var result models.BankListResponse
	err := s.client.GETWithQuery(ctx, "/banks", query, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get banks: %w", err)
	}

	return &result, nil
}

// GetAllBanks retrieves all banks with pagination handling
func (s *MiscService) GetAllBanks(ctx context.Context) ([]models.Bank, error) {
	var allBanks []models.Bank
	page := 1
	limit := 100

	for {
		opts := &models.BankListOptions{
			Page:  page,
			Limit: limit,
		}

		response, err := s.GetBanks(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get banks page %d: %w", page, err)
		}

		allBanks = append(allBanks, response.Results...)

		// Check if we have more pages
		if len(response.Results) < limit {
			break
		}
		page++
	}

	return allBanks, nil
}

// GetBankByID retrieves a specific bank by ID
func (s *MiscService) GetBankByID(ctx context.Context, bankID string) (*models.Bank, error) {
	var result models.Bank
	endpoint := fmt.Sprintf("/banks/%s", bankID)
	
	err := s.client.GET(ctx, endpoint, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank %s: %w", bankID, err)
	}

	return &result, nil
}

// GetActiveBanks retrieves only active banks
func (s *MiscService) GetActiveBanks(ctx context.Context) ([]models.Bank, error) {
	active := true
	opts := &models.BankListOptions{
		Active: &active,
		Limit:  1000, // Get all active banks
	}

	result, err := s.GetBanks(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get active banks: %w", err)
	}

	return result.Results, nil
}

// GetBanksByCountry retrieves banks for a specific country
func (s *MiscService) GetBanksByCountry(ctx context.Context, country string) ([]models.Bank, error) {
	opts := &models.BankListOptions{
		Country: country,
		Limit:   1000, // Get all banks for country
	}

	result, err := s.GetBanks(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get banks for country %s: %w", country, err)
	}

	return result.Results, nil
}

// SearchBanks searches for banks by name
func (s *MiscService) SearchBanks(ctx context.Context, searchTerm string, limit int) ([]models.Bank, error) {
	if limit <= 0 {
		limit = 50
	}

	// Get all banks and filter by name
	allBanks, err := s.GetAllBanks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search banks: %w", err)
	}

	var matchedBanks []models.Bank
	count := 0
	for _, bank := range allBanks {
		if count >= limit {
			break
		}
		// Simple case-insensitive name search
		if len(bank.Name) >= len(searchTerm) {
			for i := 0; i <= len(bank.Name)-len(searchTerm); i++ {
				if bank.Name[i:i+len(searchTerm)] == searchTerm {
					matchedBanks = append(matchedBanks, bank)
					count++
					break
				}
			}
		}
	}

	return matchedBanks, nil
}

// Business Types Section

// GetBusinessTypes retrieves a paginated list of business types
func (s *MiscService) GetBusinessTypes(ctx context.Context, opts *models.BusinessTypeListOptions) (*models.BusinessTypeListResponse, error) {
	query := url.Values{}
	
	if opts != nil {
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			query.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Country != "" {
			query.Set("country", opts.Country)
		}
		if opts.Active != nil {
			query.Set("active", strconv.FormatBool(*opts.Active))
		}
	}

	var result models.BusinessTypeListResponse
	err := s.client.GETWithQuery(ctx, "/business", query, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get business types: %w", err)
	}

	return &result, nil
}

// GetAllBusinessTypes retrieves all business types with pagination handling
func (s *MiscService) GetAllBusinessTypes(ctx context.Context) ([]models.BusinessType, error) {
	var allBusinessTypes []models.BusinessType
	page := 1
	limit := 100

	for {
		opts := &models.BusinessTypeListOptions{
			Page:  page,
			Limit: limit,
		}

		response, err := s.GetBusinessTypes(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get business types page %d: %w", page, err)
		}

		allBusinessTypes = append(allBusinessTypes, response.Results...)

		// Check if we have more pages
		if len(response.Results) < limit {
			break
		}
		page++
	}

	return allBusinessTypes, nil
}

// GetBusinessTypeByID retrieves a specific business type by ID
func (s *MiscService) GetBusinessTypeByID(ctx context.Context, businessTypeID string) (*models.BusinessType, error) {
	var result models.BusinessType
	endpoint := fmt.Sprintf("/business/%s", businessTypeID)
	
	err := s.client.GET(ctx, endpoint, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get business type %s: %w", businessTypeID, err)
	}

	return &result, nil
}

// GetBusinessTypesByCountry retrieves business types for a specific country
func (s *MiscService) GetBusinessTypesByCountry(ctx context.Context, country string) ([]models.BusinessType, error) {
	opts := &models.BusinessTypeListOptions{
		Country: country,
		Limit:   1000, // Get all business types for country
	}

	result, err := s.GetBusinessTypes(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get business types for country %s: %w", country, err)
	}

	return result.Results, nil
}

// GetActiveBusinessTypes retrieves only active business types
func (s *MiscService) GetActiveBusinessTypes(ctx context.Context) ([]models.BusinessType, error) {
	active := true
	opts := &models.BusinessTypeListOptions{
		Active: &active,
		Limit:  1000, // Get all active business types
	}

	result, err := s.GetBusinessTypes(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get active business types: %w", err)
	}

	return result.Results, nil
}

// SearchBusinessTypes searches for business types by name
func (s *MiscService) SearchBusinessTypes(ctx context.Context, searchTerm string, limit int) ([]models.BusinessType, error) {
	if limit <= 0 {
		limit = 50
	}

	// Get all business types and filter by name
	allTypes, err := s.GetAllBusinessTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search business types: %w", err)
	}

	var matchedTypes []models.BusinessType
	count := 0
	for _, businessType := range allTypes {
		if count >= limit {
			break
		}
		// Simple case-insensitive name search
		if len(businessType.Name) >= len(searchTerm) {
			for i := 0; i <= len(businessType.Name)-len(searchTerm); i++ {
				if businessType.Name[i:i+len(searchTerm)] == searchTerm {
					matchedTypes = append(matchedTypes, businessType)
					count++
					break
				}
			}
		}
	}

	return matchedTypes, nil
}