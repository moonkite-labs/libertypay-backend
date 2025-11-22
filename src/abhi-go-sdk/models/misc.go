package models

// Bank represents a bank entity
type Bank struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Country  string `json:"country"`
	IsActive bool   `json:"isActive"`
	SwiftCode string `json:"swiftCode,omitempty"`
	BankType  string `json:"bankType,omitempty"`
}

// BankListOptions represents query options for listing banks
type BankListOptions struct {
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Country string `json:"country,omitempty"`
	Active  *bool  `json:"active,omitempty"`
}

// BankListResponse represents the response for bank list
type BankListResponse struct {
	Total   int    `json:"total"`
	Results []Bank `json:"results"`
}

// BusinessType represents a business type entity  
type BusinessType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Country     string `json:"country"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"isActive,omitempty"`
}

// BusinessTypeListOptions represents query options for listing business types
type BusinessTypeListOptions struct {
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Country string `json:"country,omitempty"`
	Active  *bool  `json:"active,omitempty"`
}

// BusinessTypeListResponse represents the response for business type list
type BusinessTypeListResponse struct {
	Total   int            `json:"total"`
	Results []BusinessType `json:"results"`
}