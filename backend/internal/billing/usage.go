package billing

import (
    "time"
    "github.com/oFuterman/light-house/internal/models"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

// GetOrCreateMonthlyUsage finds or creates the usage record for the current month
// Uses upsert to handle race conditions safely
func GetOrCreateMonthlyUsage(db *gorm.DB, orgID uint) (*models.MonthlyUsage, error) {
    now := time.Now().UTC()
    year := now.Year()
    month := int(now.Month())
    usage := &models.MonthlyUsage{
        OrgID: orgID,
        Year:  year,
        Month: month,
    }
    // Upsert: create if not exists, otherwise return existing
    err := db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "org_id"}, {Name: "year"}, {Name: "month"}},
        DoNothing: true,
    }).Create(usage).Error
    if err != nil {
        return nil, err
    }
    // Fetch the actual record (whether just created or existing)
    err = db.Where("org_id = ? AND year = ? AND month = ?", orgID, year, month).First(usage).Error
    if err != nil {
        return nil, err
    }
    return usage, nil
}

// GetUsageSnapshot builds a complete usage snapshot for an org
// This aggregates current resource counts from the database
func GetUsageSnapshot(db *gorm.DB, orgID uint) (models.UsageSnapshot, error) {
    snapshot := models.UsageSnapshot{}
    // Get current check count (active, non-deleted)
    var checkCount int64
    if err := db.Model(&models.Check{}).
        Where("org_id = ? AND deleted_at IS NULL", orgID).
        Count(&checkCount).Error; err != nil {
        return snapshot, err
    }
    snapshot.CheckCount = int(checkCount)
    // Get current API key count
    var apiKeyCount int64
    if err := db.Model(&models.APIKey{}).
        Where("org_id = ?", orgID).
        Count(&apiKeyCount).Error; err != nil {
        return snapshot, err
    }
    snapshot.APIKeyCount = int(apiKeyCount)
    // Get monthly usage for log volume and AI calls
    monthlyUsage, err := GetOrCreateMonthlyUsage(db, orgID)
    if err != nil {
        return snapshot, err
    }
    snapshot.LogVolumeBytes = monthlyUsage.LogVolumeBytes
    snapshot.AILevel1Calls = monthlyUsage.AILevel1Calls
    snapshot.AILevel2Calls = monthlyUsage.AILevel2Calls
    snapshot.AILevel3Calls = monthlyUsage.AILevel3Calls
    // Status pages would be counted here when implemented
    snapshot.StatusPageCount = monthlyUsage.StatusPageCount
    return snapshot, nil
}

// IncrementLogVolume atomically adds bytes to the monthly log volume
// Returns the new total and any error
func IncrementLogVolume(db *gorm.DB, orgID uint, bytes int64) (int64, error) {
    usage, err := GetOrCreateMonthlyUsage(db, orgID)
    if err != nil {
        return 0, err
    }
    // Atomic increment
    err = db.Model(usage).
        Update("log_volume_bytes", gorm.Expr("log_volume_bytes + ?", bytes)).Error
    if err != nil {
        return 0, err
    }
    // Fetch updated value
    err = db.First(usage, usage.ID).Error
    if err != nil {
        return 0, err
    }
    return usage.LogVolumeBytes, nil
}

// IncrementAICalls atomically increments AI usage for a specific level
func IncrementAICalls(db *gorm.DB, orgID uint, level int) error {
    usage, err := GetOrCreateMonthlyUsage(db, orgID)
    if err != nil {
        return err
    }
    var column string
    switch level {
    case 1:
        column = "ai_level1_calls"
    case 2:
        column = "ai_level2_calls"
    case 3:
        column = "ai_level3_calls"
    default:
        return nil
    }
    return db.Model(usage).
        Update(column, gorm.Expr(column+" + 1")).Error
}

// SyncResourceCounts updates the cached resource counts in monthly usage
// Call this after creating/deleting checks, API keys, status pages
func SyncResourceCounts(db *gorm.DB, orgID uint) error {
    usage, err := GetOrCreateMonthlyUsage(db, orgID)
    if err != nil {
        return err
    }
    // Count current checks
    var checkCount int64
    if err := db.Model(&models.Check{}).
        Where("org_id = ? AND deleted_at IS NULL", orgID).
        Count(&checkCount).Error; err != nil {
        return err
    }
    // Count current API keys
    var apiKeyCount int64
    if err := db.Model(&models.APIKey{}).
        Where("org_id = ?", orgID).
        Count(&apiKeyCount).Error; err != nil {
        return err
    }
    // Update usage record
    return db.Model(usage).Updates(map[string]interface{}{
        "check_count":   checkCount,
        "api_key_count": apiKeyCount,
    }).Error
}

// GetCurrentCheckCount returns the current number of active checks for an org
func GetCurrentCheckCount(db *gorm.DB, orgID uint) (int, error) {
    var count int64
    err := db.Model(&models.Check{}).
        Where("org_id = ? AND deleted_at IS NULL", orgID).
        Count(&count).Error
    return int(count), err
}

// GetCurrentAPIKeyCount returns the current number of API keys for an org
func GetCurrentAPIKeyCount(db *gorm.DB, orgID uint) (int, error) {
    var count int64
    err := db.Model(&models.APIKey{}).
        Where("org_id = ?", orgID).
        Count(&count).Error
    return int(count), err
}

// GetCurrentLogVolume returns the current month's log volume for an org
func GetCurrentLogVolume(db *gorm.DB, orgID uint) (int64, error) {
    usage, err := GetOrCreateMonthlyUsage(db, orgID)
    if err != nil {
        return 0, err
    }
    return usage.LogVolumeBytes, nil
}
