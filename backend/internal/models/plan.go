package models

// Plan represents the subscription tier for an organization
type Plan string

const (
    PlanFree     Plan = "free"
    PlanIndiePro Plan = "indie_pro"
    PlanTeam     Plan = "team"
    PlanAgency   Plan = "agency"
)

// IsValid checks if the plan is a valid value
func (p Plan) IsValid() bool {
    switch p {
    case PlanFree, PlanIndiePro, PlanTeam, PlanAgency:
        return true
    }
    return false
}

// IsPaid returns true if this is a paid plan
func (p Plan) IsPaid() bool {
    return p != PlanFree
}

// PlanConfig defines the entitlements for each plan tier
type PlanConfig struct {
    Name                    string
    MaxChecks               int
    CheckIntervalMinSeconds int
    LogRetentionDays        int
    LogVolumeBytesPerMonth  int64
    MaxStatusPages          int
    MaxAPIKeys              int
    AuditLogRetentionDays   int
    AILevel1Limit           int // -1 = unlimited
    AILevel2Limit           int // -1 = unlimited
    AILevel3Limit           int // -1 = unlimited
    StripePriceID           string
    MonthlyPriceCents       int
}

// PlanConfigs maps plan tiers to their configuration
// These values align with the monetization strategy
var PlanConfigs = map[Plan]PlanConfig{
    PlanFree: {
        Name:                    "Free",
        MaxChecks:               10,
        CheckIntervalMinSeconds: 300, // 5 minutes
        LogRetentionDays:        7,
        LogVolumeBytesPerMonth:  500 * 1024 * 1024, // 500 MB
        MaxStatusPages:          0,
        MaxAPIKeys:              2,
        AuditLogRetentionDays:   0,
        AILevel1Limit:           1,  // 1 per day
        AILevel2Limit:           0,
        AILevel3Limit:           0,
        StripePriceID:           "", // No Stripe subscription for free
        MonthlyPriceCents:       0,
    },
    PlanIndiePro: {
        Name:                    "Indie Pro",
        MaxChecks:               25,
        CheckIntervalMinSeconds: 60, // 1 minute
        LogRetentionDays:        30,
        LogVolumeBytesPerMonth:  5 * 1024 * 1024 * 1024, // 5 GB
        MaxStatusPages:          1,
        MaxAPIKeys:              10,
        AuditLogRetentionDays:   7,
        AILevel1Limit:           -1, // unlimited
        AILevel2Limit:           30,
        AILevel3Limit:           0,
        StripePriceID:           "", // Set via environment variable
        MonthlyPriceCents:       1900,
    },
    PlanTeam: {
        Name:                    "Team",
        MaxChecks:               75,
        CheckIntervalMinSeconds: 30, // 30 seconds
        LogRetentionDays:        90,
        LogVolumeBytesPerMonth:  20 * 1024 * 1024 * 1024, // 20 GB
        MaxStatusPages:          3,
        MaxAPIKeys:              25,
        AuditLogRetentionDays:   30,
        AILevel1Limit:           -1, // unlimited
        AILevel2Limit:           -1, // unlimited
        AILevel3Limit:           30,
        StripePriceID:           "", // Set via environment variable
        MonthlyPriceCents:       4900,
    },
    PlanAgency: {
        Name:                    "Agency",
        MaxChecks:               250,
        CheckIntervalMinSeconds: 30, // 30 seconds
        LogRetentionDays:        180,
        LogVolumeBytesPerMonth:  50 * 1024 * 1024 * 1024, // 50 GB
        MaxStatusPages:          -1,                      // unlimited
        MaxAPIKeys:              -1,                      // unlimited
        AuditLogRetentionDays:   365,
        AILevel1Limit:           -1, // unlimited
        AILevel2Limit:           -1, // unlimited
        AILevel3Limit:           -1, // unlimited
        StripePriceID:           "", // Set via environment variable
        MonthlyPriceCents:       14900,
    },
}

// GetPlanConfig returns the configuration for a plan
func GetPlanConfig(plan Plan) PlanConfig {
    if config, ok := PlanConfigs[plan]; ok {
        return config
    }
    // Default to free if unknown plan
    return PlanConfigs[PlanFree]
}

// GetOrgPlanConfig is a convenience function that extracts the plan from an org
// and returns the corresponding config. Logs a warning for invalid plans.
func GetOrgPlanConfig(org *Organization) PlanConfig {
    if org == nil {
        return PlanConfigs[PlanFree]
    }
    if !org.Plan.IsValid() {
        // Invalid plan - default to free (logging should be done by caller)
        return PlanConfigs[PlanFree]
    }
    return GetPlanConfig(org.Plan)
}

// GetPlanByStripePriceID finds the plan for a Stripe price ID
func GetPlanByStripePriceID(priceID string) (Plan, bool) {
    for plan, config := range PlanConfigs {
        if config.StripePriceID == priceID {
            return plan, true
        }
    }
    return PlanFree, false
}
