package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oFuterman/light-house/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// JWTSecret is set during router initialization
var JWTSecret string

// AuthRequired validates JWT tokens and sets user context
// Accepts token from cookie (browser) or Authorization header (API clients)
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string

		// Try cookie first (browser clients)
		tokenString = c.Cookies("token")

		// Fall back to Authorization header (API clients)
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "missing authorization",
				})
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid authorization header format",
				})
			}
			tokenString = parts[1]
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token signing method")
			}
			return []byte(JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		// Get user_id and org_id from claims
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid user_id in token",
			})
		}

		orgIDFloat, ok := claims["org_id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid org_id in token",
			})
		}

		// Get role from claims (default to member for backwards compatibility)
		role := models.RoleMember
		if roleStr, ok := claims["role"].(string); ok {
			role = models.Role(roleStr)
		}

		// Set user info in context for handlers to use
		c.Locals("userID", uint(userIDFloat))
		c.Locals("orgID", uint(orgIDFloat))
		c.Locals("role", role)

		return c.Next()
	}
}

// AuthRequiredWithDB validates JWT and fetches fresh user data from DB
// Use this when you need up-to-date role information
func AuthRequiredWithDB(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string

		// Try cookie first (browser clients)
		tokenString = c.Cookies("token")

		// Fall back to Authorization header (API clients)
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "missing authorization",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid authorization header format",
				})
			}
			tokenString = parts[1]
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token signing method")
			}
			return []byte(JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid user_id in token",
			})
		}

		// Fetch user from database for fresh role info
		var user models.User
		if err := db.First(&user, uint(userIDFloat)).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "user not found",
			})
		}

		c.Locals("userID", user.ID)
		c.Locals("orgID", user.OrgID)
		c.Locals("role", user.Role)
		c.Locals("user", &user)

		return c.Next()
	}
}

// RequireRole creates middleware that checks if the user has the required role
func RequireRole(roles ...models.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(models.Role)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "role not available",
			})
		}

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "insufficient permissions",
		})
	}
}

// RequireAdmin is a shortcut for RequireRole(owner, admin)
func RequireAdmin() fiber.Handler {
	return RequireRole(models.RoleOwner, models.RoleAdmin)
}

// RequireOwner is a shortcut for RequireRole(owner)
func RequireOwner() fiber.Handler {
	return RequireRole(models.RoleOwner)
}

// APIKeyAuth validates API keys for log ingestion
func APIKeyAuth(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing X-API-Key header",
			})
		}

		// Find API key by checking hash
		var keys []models.APIKey
		if err := db.Find(&keys).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to validate API key",
			})
		}

		var matchedKey *models.APIKey
		for _, key := range keys {
			if err := bcrypt.CompareHashAndPassword([]byte(key.KeyHash), []byte(apiKey)); err == nil {
				matchedKey = &key
				break
			}
		}

		if matchedKey == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid API key",
			})
		}

		// Update last used timestamp
		now := time.Now()
		db.Model(matchedKey).Update("last_used_at", now)

		// Set context
		c.Locals("orgID", matchedKey.OrgID)
		c.Locals("apiKey", matchedKey)
		c.Locals("apiKeyScopes", matchedKey.Scopes)

		return c.Next()
	}
}

// APIKeyAuthWithScope validates API keys and checks for required scopes
func APIKeyAuthWithScope(db *gorm.DB, requiredScopes ...models.APIKeyScope) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing X-API-Key header",
			})
		}

		// Find API key by checking hash
		var keys []models.APIKey
		if err := db.Find(&keys).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to validate API key",
			})
		}

		var matchedKey *models.APIKey
		for _, key := range keys {
			if err := bcrypt.CompareHashAndPassword([]byte(key.KeyHash), []byte(apiKey)); err == nil {
				matchedKey = &key
				break
			}
		}

		if matchedKey == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid API key",
			})
		}

		// Check scopes
		if len(requiredScopes) > 0 && !matchedKey.HasAnyScope(requiredScopes...) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "API key lacks required scope",
			})
		}

		// Update last used timestamp
		now := time.Now()
		db.Model(matchedKey).Update("last_used_at", now)

		// Set context
		c.Locals("orgID", matchedKey.OrgID)
		c.Locals("apiKey", matchedKey)
		c.Locals("apiKeyScopes", matchedKey.Scopes)

		return c.Next()
	}
}
