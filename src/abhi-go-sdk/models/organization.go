package models

import "time"

// Organization represents an organization entity
type Organization struct {
	ID                     string       `json:"id,omitempty"`
	DeletedDate            *time.Time   `json:"deletedDate,omitempty"`
	Name                   string       `json:"name" validate:"required"`
	OrganizationNumber     int          `json:"organizationNumber,omitempty"`
	Industry               string       `json:"industry" validate:"required"`
	ManagementAlias        string       `json:"managementAlias" validate:"required,min=4,max=100,lowercase"`
	OrganizationType       string       `json:"organizationType,omitempty"`
	Active                 bool         `json:"active,omitempty"`
	CreditLimit            float64      `json:"creditLimit" validate:"required,gt=0"`
	Address                string       `json:"address" validate:"required"`
	City                   string       `json:"city" validate:"required"`
	Phone                  string       `json:"phone,omitempty"`
	Email                  string       `json:"email,omitempty,email"`
	PayrollStartDay        int          `json:"payrollStartDay,omitempty" validate:"omitempty,min=1,max=31"`
	BusinessTypeID         string       `json:"businessTypeId" validate:"required,uuid4"`
	BusinessType           BusinessType `json:"businessType,omitempty"`
	ParentOrganizationID   string       `json:"parentOrganizationId,omitempty"`
	ParentOrganizations    *ParentOrg   `json:"parentOrganizations,omitempty"`
	CreatedAt              time.Time    `json:"createdAt,omitempty"`
	UpdatedAt              time.Time    `json:"updatedAt,omitempty"`
}

// ParentOrg represents parent organization relationship
type ParentOrg struct {
	ParentOrganizationID string `json:"parentOrganizationId"`
}


// CreateOrganizationRequest represents a request to create a new organization
type CreateOrganizationRequest struct {
	Name            string  `json:"name" validate:"required"`
	Industry        string  `json:"industry" validate:"required"`
	BusinessTypeID  string  `json:"businessTypeId" validate:"required,uuid4"`
	Address         string  `json:"address" validate:"required"`
	City            string  `json:"city" validate:"required"`
	ManagementAlias string  `json:"managementAlias" validate:"required,min=4,max=100"`
	CreditLimit     float64 `json:"creditLimit" validate:"required,gt=0"`
	Phone           string  `json:"phone,omitempty"`
	Email           string  `json:"email,omitempty,email"`
	PayrollStartDay int     `json:"payrollStartDay,omitempty" validate:"omitempty,min=1,max=31"`
}

// OrganizationListOptions represents query options for listing organizations
type OrganizationListOptions struct {
	Page         int    `json:"page,omitempty"`
	Limit        int    `json:"limit,omitempty"`
	From         string `json:"from,omitempty"`         // Start date filter (ISO format)
	To           string `json:"to,omitempty"`           // End date filter (ISO format)
	ShowInactive bool   `json:"showInactive,omitempty"` // Include inactive organizations
	Column       string `json:"column,omitempty"`       // Sort column: "organizations.createdAt", "organizations.name"
	Order        string `json:"order,omitempty"`        // Sort order: "ASC", "DESC"
}

// OrganizationListResponse represents the response for organization list
type OrganizationListResponse struct {
	Total   int            `json:"total"`
	Results []Organization `json:"results"`
}

// CreateOrganizationResponse represents the response when creating an organization
type CreateOrganizationResponse struct {
	Message string                     `json:"message"`
	Data    OrganizationCreationResult `json:"data"`
}

// OrganizationCreationResult contains the created organization details and users
type OrganizationCreationResult struct {
	Users          OrganizationUsers `json:"users"`
	OrganizationID string            `json:"organizationId"`
}

// OrganizationUsers contains the automatically created users for the organization
type OrganizationUsers struct {
	Admin        OrganizationUser `json:"admin"`
	Operator     OrganizationUser `json:"operator"`
	Support      OrganizationUser `json:"support,omitempty"`
}

// OrganizationUser represents a user created for an organization
type OrganizationUser struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	MFASecret string `json:"MFAsecret,omitempty"`
}

