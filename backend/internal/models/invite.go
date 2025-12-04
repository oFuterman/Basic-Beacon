package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

// InviteStatus represents the state of an invitation
type InviteStatus string

const (
	InviteStatusPending  InviteStatus = "pending"
	InviteStatusAccepted InviteStatus = "accepted"
	InviteStatusExpired  InviteStatus = "expired"
	InviteStatusRevoked  InviteStatus = "revoked"
)

// Invite represents an invitation to join an organization
type Invite struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OrgID       uint         `gorm:"not null;index" json:"org_id"`
	Email       string       `gorm:"not null;size:255;index" json:"email"`
	Role        Role         `gorm:"not null;size:20;default:'member'" json:"role"`
	Token       string       `gorm:"not null;uniqueIndex;size:64" json:"-"`
	Status      InviteStatus `gorm:"not null;size:20;default:'pending'" json:"status"`
	ExpiresAt   time.Time    `gorm:"not null" json:"expires_at"`
	InvitedByID uint         `gorm:"not null" json:"invited_by_id"`
	AcceptedAt  *time.Time   `json:"accepted_at,omitempty"`

	// Relations
	Organization Organization `gorm:"foreignKey:OrgID" json:"organization,omitempty"`
	InvitedBy    User         `gorm:"foreignKey:InvitedByID" json:"invited_by,omitempty"`
}

// GenerateToken creates a cryptographically secure random token
func GenerateInviteToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsExpired checks if the invite has passed its expiration time
func (i *Invite) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// IsUsable checks if the invite can still be accepted
func (i *Invite) IsUsable() bool {
	return i.Status == InviteStatusPending && !i.IsExpired()
}

// DefaultInviteExpiration returns the default duration for invite validity
func DefaultInviteExpiration() time.Duration {
	return 7 * 24 * time.Hour // 7 days
}
