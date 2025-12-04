package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// APIKeyScope represents a permission scope for API keys
type APIKeyScope string

const (
	// Ingestion scopes
	ScopeLogsWrite   APIKeyScope = "logs:write"
	ScopeTracesWrite APIKeyScope = "traces:write"

	// Read scopes
	ScopeLogsRead    APIKeyScope = "logs:read"
	ScopeTracesRead  APIKeyScope = "traces:read"
	ScopeChecksRead  APIKeyScope = "checks:read"
	ScopeAlertsRead  APIKeyScope = "alerts:read"

	// Write scopes
	ScopeChecksWrite APIKeyScope = "checks:write"

	// Full access
	ScopeAll APIKeyScope = "*"
)

// AllScopes returns all available scopes
func AllScopes() []APIKeyScope {
	return []APIKeyScope{
		ScopeLogsWrite,
		ScopeTracesWrite,
		ScopeLogsRead,
		ScopeTracesRead,
		ScopeChecksRead,
		ScopeAlertsRead,
		ScopeChecksWrite,
		ScopeAll,
	}
}

// IsValidScope checks if a scope string is valid
func IsValidScope(scope string) bool {
	for _, s := range AllScopes() {
		if string(s) == scope {
			return true
		}
	}
	return false
}

type APIKey struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OrgID       uint           `gorm:"not null;index" json:"org_id"`
	Name        string         `gorm:"not null;size:255" json:"name"`
	KeyHash     string         `gorm:"not null;uniqueIndex" json:"-"`
	Prefix      string         `gorm:"not null;size:12" json:"prefix"` // First 8 chars for identification
	Scopes      pq.StringArray `gorm:"type:text[];not null;default:'{}'" json:"scopes"`
	LastUsedAt  *time.Time     `json:"last_used_at,omitempty"`
	CreatedByID *uint          `json:"created_by_id,omitempty"`

	// Relations
	Organization Organization `gorm:"foreignKey:OrgID" json:"organization,omitempty"`
	CreatedBy    *User        `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
}

// HasScope checks if the API key has a specific scope
func (k *APIKey) HasScope(scope APIKeyScope) bool {
	for _, s := range k.Scopes {
		if s == string(ScopeAll) || s == string(scope) {
			return true
		}
	}
	return false
}

// HasAnyScope checks if the API key has any of the specified scopes
func (k *APIKey) HasAnyScope(scopes ...APIKeyScope) bool {
	for _, scope := range scopes {
		if k.HasScope(scope) {
			return true
		}
	}
	return false
}
