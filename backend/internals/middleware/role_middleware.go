package middleware

import (
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

// RequireRole creates a Gin middleware that checks for a specific Clerk role.
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the claims that your ClerkAuth middleware already validated.
		claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
		if !ok || claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Check role in organization role (e.g. "admin", "org:admin") or custom claims.
		hasRole := claims.HasRole(role) ||
			claims.ActiveOrganizationRole == role ||
			claims.ActiveOrganizationRole == "org:"+role ||
			(strings.HasPrefix(role, "org:") && claims.ActiveOrganizationRole == strings.TrimPrefix(role, "org:"))

		if !hasRole && claims.Custom != nil {
			if m, ok := claims.Custom.(map[string]any); ok {
				if r, ok := m["role"].(string); ok && (r == role || r == "org:"+role) {
					hasRole = true
				}
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}

		c.Next()
	}
}
