package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORSMiddleware("http://localhost:5173, https://app.example.com/"))
	r.GET("/api/products", func(c *gin.Context) { c.Status(http.StatusOK) })

	for _, tc := range []struct {
		name, method, origin string
		status               int
		allowed              bool
	}{
		{"allowed preflight", http.MethodOptions, "http://localhost:5173", http.StatusNoContent, true},
		{"allowed request", http.MethodGet, "https://app.example.com", http.StatusOK, true},
		{"rejected preflight", http.MethodOptions, "https://evil.example.com", http.StatusForbidden, false},
		{"rejected request", http.MethodGet, "https://evil.example.com", http.StatusOK, false},
		{"no origin", http.MethodGet, "", http.StatusOK, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/api/products", nil)
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			response := httptest.NewRecorder()
			r.ServeHTTP(response, req)
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d", response.Code, tc.status)
			}
			wantOrigin := ""
			if tc.allowed {
				wantOrigin = tc.origin
			}
			if got := response.Header().Get("Access-Control-Allow-Origin"); got != wantOrigin {
				t.Errorf("allow origin = %q, want %q", got, wantOrigin)
			}
			if got := response.Header().Get("Access-Control-Allow-Credentials"); got != "" {
				t.Errorf("unexpected credentials header: %q", got)
			}
			if got := response.Header().Get("Vary"); got != "Origin" {
				t.Errorf("Vary = %q, want Origin", got)
			}
			if tc.name == "allowed preflight" && response.Header().Get("Access-Control-Allow-Headers") != "Authorization, Content-Type" {
				t.Error("Authorization and Content-Type must be allowed")
			}
		})
	}
}
