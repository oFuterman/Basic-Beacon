package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents the permission level of a user within an organization
type Role string

const (
	RoleOwner  Role = "owner"  // Full access, can delete org, manage billing
	RoleAdmin  Role = "admin"  // Can manage members, settings, and all resources
	RoleMember Role = "member" // Can view and create resources, but limited management
)

// IsValid checks if the role is a valid value
func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleAdmin, RoleMember:
		return true
	}
	return false
}

// CanManageMembers returns true if the role can invite/remove members
func (r Role) CanManageMembers() bool {
	return r == RoleOwner || r == RoleAdmin
}

// CanManageSettings returns true if the role can modify org settings
func (r Role) CanManageSettings() bool {
	return r == RoleOwner || r == RoleAdmin
}

// CanDeleteOrg returns true if the role can delete the organization
func (r Role) CanDeleteOrg() bool {
	return r == RoleOwner
}

// CanManageAPIKeys returns true if the role can create/delete API keys
func (r Role) CanManageAPIKeys() bool {
	return r == RoleOwner || r == RoleAdmin
}

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email        string `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	OrgID        uint   `gorm:"not null;index" json:"org_id"`
	Role         Role   `gorm:"not null;size:20;default:'member'" json:"role"`

	// Relations
	Organization Organization `gorm:"foreignKey:OrgID" json:"organization,omitempty"`
}
