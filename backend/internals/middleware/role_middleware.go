package middleware

import (
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
)

// RequireRole allows the request through only if the Clerk session claims
// carry the given role, either as an organization role ("org:admin") or as
// a custom "role" claim.
func RequireRole(role string) gin.HandlerFunc {
	want := normalizeRole(role)

	return func(c *gin.Context) {
		claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
		if !ok || claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if !hasRole(claims, want) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}

		c.Next()
	}
}

func hasRole(claims *clerk.SessionClaims, want string) bool {
	if normalizeRole(claims.ActiveOrganizationRole) == want {
		return true
	}

	custom, ok := claims.Custom.(map[string]any)
	if !ok {
		return false
	}
	role, ok := custom["role"].(string)
	return ok && normalizeRole(role) == want
}

func normalizeRole(role string) string {
	return strings.ToLower(strings.TrimPrefix(role, "org:"))
}
