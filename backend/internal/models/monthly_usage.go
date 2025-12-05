package models

import (
    "time"
)

// MonthlyUsage tracks resource consumption per organization per month
// Used for enforcement and UI display, NOT for billing calculations
type MonthlyUsage struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`

    // Composite unique key: one record per org per month
    OrgID uint `gorm:"uniqueIndex:idx_usage_org_month;not null" json:"org_id"`
    Year  int  `gorm:"uniqueIndex:idx_usage_org_month;not null" json:"year"`
    Month int  `gorm:"uniqueIndex:idx_usage_org_month;not null" json:"month"` // 1-12

    // Resource counts (current state, not deltas)
    LogVolumeBytes  int64 `gorm:"default:0" json:"log_volume_bytes"`
    CheckCount      int   `gorm:"default:0" json:"check_count"`       // Current active checks
    StatusPageCount int   `gorm:"default:0" json:"status_page_count"` // Current status pages
    APIKeyCount     int   `gorm:"default:0" json:"api_key_count"`     // Current API keys

    // AI usage (resets monthly)
    AILevel1Calls int `gorm:"default:0" json:"ai_level1_calls"`
    AILevel2Calls int `gorm:"default:0" json:"ai_level2_calls"`
    AILevel3Calls int `gorm:"default:0" json:"ai_level3_calls"`
}

// UsageSnapshot represents a point-in-time view of an org's resource usage
// This is passed to the entitlement engine to check limits
type UsageSnapshot struct {
    CheckCount      int   `json:"check_count"`
    LogVolumeBytes  int64 `json:"log_volume_bytes"`
    StatusPageCount int   `json:"status_page_count"`
    APIKeyCount     int   `json:"api_key_count"`
    AILevel1Calls   int   `json:"ai_level1_calls"`
    AILevel2Calls   int   `json:"ai_level2_calls"`
    AILevel3Calls   int   `json:"ai_level3_calls"`
}

// ToSnapshot converts MonthlyUsage to a UsageSnapshot
func (m *MonthlyUsage) ToSnapshot() UsageSnapshot {
    return UsageSnapshot{
        CheckCount:      m.CheckCount,
        LogVolumeBytes:  m.LogVolumeBytes,
        StatusPageCount: m.StatusPageCount,
        APIKeyCount:     m.APIKeyCount,
        AILevel1Calls:   m.AILevel1Calls,
        AILevel2Calls:   m.AILevel2Calls,
        AILevel3Calls:   m.AILevel3Calls,
    }
}
