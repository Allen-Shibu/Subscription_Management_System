package middleware

import (
	"context"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gin-gonic/gin"
)

// RoleClaims is the custom session claim added in the Clerk dashboard.
type RoleClaims struct {
	Role string `json:"role"`
}

// ClerkAuth authenticates incoming requests using Clerk session tokens.
func ClerkAuth() gin.HandlerFunc {
	// Clerk's standard net/http middleware.
	mw := clerkhttp.WithHeaderAuthorization(
		clerkhttp.AuthorizationFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
		})),
		// Without this, custom claims (role) are dropped from the verified token.
		clerkhttp.CustomClaimsConstructor(func(context.Context) any {
			return &RoleClaims{}
		}),
	)

	return func(c *gin.Context) {
		nextCalled := false

		// Clerk writes the verified SessionClaims into r.Context().
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			c.Request = r // carry the enriched request context into Gin
		})

		mw(next).ServeHTTP(c.Writer, c.Request)

		// Clerk rejected the token and did not call next.
		if !nextCalled {
			c.Abort()
			return
		}

		// WithHeaderAuthorization allows requests with no token through.
		// Explicitly reject them here.
		claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
		if !ok || claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		c.Next()
	}
}

// GetClaims retrieves validated session claims from Gin context.
func GetClaims(c *gin.Context) (*clerk.SessionClaims, bool) {
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())
	return claims, ok && claims != nil
}

// GetUserID retrieves the authenticated user's ID from claims.
func GetUserID(c *gin.Context) (string, bool) {
	claims, ok := GetClaims(c)
	if !ok || claims == nil {
		return "", false
	}
	return claims.Subject, claims.Subject != ""
}
