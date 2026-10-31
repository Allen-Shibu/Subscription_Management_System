package main

import (
	"log"
	"net/http"
	"os"

	"recur_backend/internals/middleware"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Apply CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Public health check endpoints
	healthHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "recur_backend",
		})
	}
	r.GET("/health", healthHandler)

	// Protected API routes
	api := r.Group("/api")
	api.Use(middleware.ClerkAuth())
	{
		// Profile endpoint returning current authenticated user ID & claims
		api.GET("/me", func(c *gin.Context) {
			userID, _ := middleware.GetUserID(c)
			claims, _ := middleware.GetClaims(c)

			var role, orgID string
			var custom any
			if claims != nil {
				role = claims.ActiveOrganizationRole
				orgID = claims.ActiveOrganizationID
				custom = claims.Custom
			}

			c.JSON(http.StatusOK, gin.H{
				"userId": userID,
				"role":   role,
				"orgId":  orgID,
				"custom": custom,
			})
		})

		// Admin-only sub-group example
		admin := api.Group("/admin")
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status":  "ok",
					"message": "welcome admin",
				})
			})
		}
	}

	return r
}

func main() {
	// Load .env file if available (ignore error if missing)
	_ = godotenv.Load()

	secretKey := os.Getenv("CLERK_SECRET_KEY")
	if secretKey != "" {
		clerk.SetKey(secretKey)
		log.Println("[INFO] Clerk secret key configured.")
	} else {
		log.Println("[WARN] CLERK_SECRET_KEY is not set. Authenticated endpoints will reject requests.")
	}

	r := setupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("[INFO] Server starting on port :%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
