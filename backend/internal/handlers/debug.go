package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/oFuterman/light-house/internal/billing"
    "github.com/oFuterman/light-house/internal/models"
    "gorm.io/gorm"
)

// DebugEntitlementsResponse is the response for the debug entitlements endpoint
type DebugEntitlementsResponse struct {
    OrgID        uint                     `json:"org_id"`
    Plan         string                   `json:"plan"`
    PlanConfig   models.PlanConfig        `json:"plan_config"`
    Usage        models.UsageSnapshot     `json:"usage"`
    Entitlements billing.EntitlementResult `json:"entitlements"`
}

// GetDebugEntitlements returns the current org's plan, config, usage, and entitlements
// This is a debug-only endpoint for development and testing
// GET /api/v1/debug/entitlements
func GetDebugEntitlements(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Only allow in development environment
        if Environment != "development" && Environment != "test" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "debug endpoints are only available in development mode",
            })
        }
        orgID := c.Locals("orgID").(uint)
        // Load org to get plan
        var org models.Organization
        if err := db.First(&org, orgID).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "failed to load organization",
            })
        }
        // Get plan config
        planConfig := models.GetOrgPlanConfig(&org)
        // Get usage snapshot
        usage, err := billing.GetUsageSnapshot(db, orgID)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "failed to get usage snapshot",
            })
        }
        // Evaluate entitlements
        entitlements := billing.CheckEntitlements(org.Plan, usage)
        return c.JSON(DebugEntitlementsResponse{
            OrgID:        orgID,
            Plan:         string(org.Plan),
            PlanConfig:   planConfig,
            Usage:        usage,
            Entitlements: entitlements,
        })
    }
}
