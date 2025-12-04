package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/oFuterman/light-house/internal/models"
	"gorm.io/gorm"
)

// MemberResponse is the response DTO for organization members
type MemberResponse struct {
	ID        uint        `json:"id"`
	Email     string      `json:"email"`
	Role      models.Role `json:"role"`
	CreatedAt string      `json:"created_at"`
}

// UpdateMemberRoleRequest is the request body for updating a member's role
type UpdateMemberRoleRequest struct {
	Role models.Role `json:"role"`
}

// ListMembers returns all members of the organization
func ListMembers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := c.Locals("orgID").(uint)

		var users []models.User
		if err := db.Where("org_id = ?", orgID).Order("created_at ASC").Find(&users).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to fetch members",
			})
		}

		// Convert to response DTOs
		members := make([]MemberResponse, len(users))
		for i, user := range users {
			members[i] = MemberResponse{
				ID:        user.ID,
				Email:     user.Email,
				Role:      user.Role,
				CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			}
		}

		return c.JSON(fiber.Map{
			"members": members,
		})
	}
}

// GetMember returns a single member by ID
func GetMember(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := c.Locals("orgID").(uint)
		memberID := c.Params("id")

		var user models.User
		if err := db.Where("id = ? AND org_id = ?", memberID, orgID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "member not found",
			})
		}

		return c.JSON(MemberResponse{
			ID:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
}

// UpdateMemberRole updates a member's role
func UpdateMemberRole(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := c.Locals("orgID").(uint)
		currentUserID := c.Locals("userID").(uint)
		userRole := c.Locals("role").(models.Role)

		if !userRole.CanManageMembers() {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "insufficient permissions",
			})
		}

		memberID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid member ID",
			})
		}

		var req UpdateMemberRoleRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		if !req.Role.IsValid() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid role",
			})
		}

		// Find the member
		var member models.User
		if err := db.Where("id = ? AND org_id = ?", memberID, orgID).First(&member).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "member not found",
			})
		}

		// Prevent self-demotion for owners
		if uint(memberID) == currentUserID && member.Role == models.RoleOwner && req.Role != models.RoleOwner {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "owners cannot demote themselves",
			})
		}

		// Only owners can change roles to/from owner/admin
		if userRole != models.RoleOwner {
			if member.Role == models.RoleOwner || member.Role == models.RoleAdmin {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "only owners can modify admin/owner roles",
				})
			}
			if req.Role == models.RoleOwner || req.Role == models.RoleAdmin {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "only owners can assign admin/owner roles",
				})
			}
		}

		// Ensure there's always at least one owner
		if member.Role == models.RoleOwner && req.Role != models.RoleOwner {
			var ownerCount int64
			db.Model(&models.User{}).Where("org_id = ? AND role = ?", orgID, models.RoleOwner).Count(&ownerCount)
			if ownerCount <= 1 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "organization must have at least one owner",
				})
			}
		}

		oldRole := member.Role
		member.Role = req.Role

		if err := db.Save(&member).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to update member role",
			})
		}

		// Log audit event
		logAuditEvent(db, orgID, &currentUserID, models.AuditActionMemberRoleChanged, "user", &member.ID, models.JSONMap{
			"email":    member.Email,
			"old_role": string(oldRole),
			"new_role": string(req.Role),
		}, c.IP(), c.Get("User-Agent"))

		return c.JSON(MemberResponse{
			ID:        member.ID,
			Email:     member.Email,
			Role:      member.Role,
			CreatedAt: member.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
}

// RemoveMember removes a member from the organization
func RemoveMember(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := c.Locals("orgID").(uint)
		currentUserID := c.Locals("userID").(uint)
		userRole := c.Locals("role").(models.Role)

		if !userRole.CanManageMembers() {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "insufficient permissions",
			})
		}

		memberID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid member ID",
			})
		}

		// Find the member
		var member models.User
		if err := db.Where("id = ? AND org_id = ?", memberID, orgID).First(&member).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "member not found",
			})
		}

		// Prevent self-removal
		if uint(memberID) == currentUserID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot remove yourself from the organization",
			})
		}

		// Only owners can remove admins/owners
		if userRole != models.RoleOwner && (member.Role == models.RoleOwner || member.Role == models.RoleAdmin) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "only owners can remove admins or owners",
			})
		}

		// Ensure there's always at least one owner
		if member.Role == models.RoleOwner {
			var ownerCount int64
			db.Model(&models.User{}).Where("org_id = ? AND role = ?", orgID, models.RoleOwner).Count(&ownerCount)
			if ownerCount <= 1 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "cannot remove the last owner of the organization",
				})
			}
		}

		// Soft delete the member
		if err := db.Delete(&member).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to remove member",
			})
		}

		// Log audit event
		logAuditEvent(db, orgID, &currentUserID, models.AuditActionMemberRemoved, "user", &member.ID, models.JSONMap{
			"email": member.Email,
			"role":  string(member.Role),
		}, c.IP(), c.Get("User-Agent"))

		return c.JSON(fiber.Map{
			"message": "member removed successfully",
		})
	}
}

// LeaveOrganization allows a member to leave the organization
func LeaveOrganization(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := c.Locals("orgID").(uint)
		userID := c.Locals("userID").(uint)
		userRole := c.Locals("role").(models.Role)

		// Check if user is the last owner
		if userRole == models.RoleOwner {
			var ownerCount int64
			db.Model(&models.User{}).Where("org_id = ? AND role = ?", orgID, models.RoleOwner).Count(&ownerCount)
			if ownerCount <= 1 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "cannot leave organization as the last owner. Transfer ownership first or delete the organization.",
				})
			}
		}

		// Get user info for audit log
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}

		// Soft delete the user
		if err := db.Delete(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to leave organization",
			})
		}

		// Log audit event
		logAuditEvent(db, orgID, &userID, models.AuditActionMemberRemoved, "user", &user.ID, models.JSONMap{
			"email":  user.Email,
			"role":   string(user.Role),
			"action": "self_leave",
		}, c.IP(), c.Get("User-Agent"))

		// Clear auth cookie
		clearAuthCookie(c)

		return c.JSON(fiber.Map{
			"message": "successfully left the organization",
		})
	}
}

// TransferOwnership transfers ownership to another member
func TransferOwnership(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := c.Locals("orgID").(uint)
		currentUserID := c.Locals("userID").(uint)
		userRole := c.Locals("role").(models.Role)

		if userRole != models.RoleOwner {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "only owners can transfer ownership",
			})
		}

		memberID, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid member ID",
			})
		}

		if uint(memberID) == currentUserID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot transfer ownership to yourself",
			})
		}

		// Find the target member
		var targetMember models.User
		if err := db.Where("id = ? AND org_id = ?", memberID, orgID).First(&targetMember).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "member not found",
			})
		}

		// Find current user
		var currentUser models.User
		if err := db.First(&currentUser, currentUserID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "current user not found",
			})
		}

		// Perform the transfer in a transaction
		err = db.Transaction(func(tx *gorm.DB) error {
			// Make target member owner
			targetMember.Role = models.RoleOwner
			if err := tx.Save(&targetMember).Error; err != nil {
				return err
			}

			// Demote current owner to admin
			currentUser.Role = models.RoleAdmin
			if err := tx.Save(&currentUser).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to transfer ownership",
			})
		}

		// Log audit event
		logAuditEvent(db, orgID, &currentUserID, models.AuditActionMemberRoleChanged, "user", &targetMember.ID, models.JSONMap{
			"action":       "ownership_transfer",
			"new_owner":    targetMember.Email,
			"former_owner": currentUser.Email,
		}, c.IP(), c.Get("User-Agent"))

		return c.JSON(fiber.Map{
			"message": "ownership transferred successfully",
		})
	}
}
