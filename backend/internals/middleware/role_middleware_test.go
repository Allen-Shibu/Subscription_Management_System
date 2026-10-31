package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
	"recur_backend/pkg/types"
)

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		claims     *clerk.SessionClaims
		wantRole   string
		wantStatus int
	}{
		{
			name:       "organization admin matches admin",
			claims:     &clerk.SessionClaims{Claims: clerk.Claims{ActiveOrganizationRole: "org:admin"}},
			wantRole:   types.RoleAdmin,
			wantStatus: http.StatusOK,
		},
		{
			name:       "role matching is case insensitive",
			claims:     &clerk.SessionClaims{Claims: clerk.Claims{ActiveOrganizationRole: "Admin"}},
			wantRole:   types.RoleAdmin,
			wantStatus: http.StatusOK,
		},
		{
			name:       "custom role claim matches",
			claims:     &clerk.SessionClaims{Custom: &RoleClaims{Role: "staff"}},
			wantRole:   types.RoleStaff,
			wantStatus: http.StatusOK,
		},
		{
			name:       "customer cannot reach admin",
			claims:     &clerk.SessionClaims{Claims: clerk.Claims{ActiveOrganizationRole: "org:customer"}},
			wantRole:   types.RoleAdmin,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "missing claims is unauthorized",
			claims:     nil,
			wantRole:   types.RoleAdmin,
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/",
				func(c *gin.Context) {
					c.Request = c.Request.WithContext(clerk.ContextWithSessionClaims(c.Request.Context(), tt.claims))
				},
				RequireRole(tt.wantRole),
				func(c *gin.Context) { c.Status(http.StatusOK) },
			)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

			if w.Code != tt.wantStatus {
				t.Fatalf("got status %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
