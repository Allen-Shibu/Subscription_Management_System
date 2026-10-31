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
	return RequireAnyRole(role)
}

// RequireAnyRole allows the request when the claims carry at least one of the
// given roles.
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	want := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		want[normalizeRole(role)] = struct{}{}
	}

	return func(c *gin.Context) {
		claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
		if !ok || claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		for role := range want {
			if hasRole(claims, role) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "insufficient permissions",
		})
	}
}

func hasRole(claims *clerk.SessionClaims, want string) bool {
	if normalizeRole(claims.ActiveOrganizationRole) == want {
		return true
	}

	custom, ok := claims.Custom.(*RoleClaims)
	return ok && normalizeRole(custom.Role) == want
}

func normalizeRole(role string) string {
	return strings.ToLower(strings.TrimPrefix(role, "org:"))
}
