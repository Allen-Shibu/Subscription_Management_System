package main

import (
	"log"
	"net/http"

	"recur_backend/config"
	"recur_backend/internals/middleware"
	"recur_backend/pkg/database"
	"recur_backend/pkg/types"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
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
		admin.Use(middleware.RequireRole(types.RoleAdmin))
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
	cfg := config.Load()

	if cfg.ClerkSecretKey == "" {
		log.Fatal("CLERK_SECRET_KEY is required")
	}
	clerk.SetKey(cfg.ClerkSecretKey)

	//setting up db connection
	pool, err := database.ConnectDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Database connection failed.")
		return
	}
	defer pool.Close()

	r := setupRouter()

	log.Printf("[INFO] Server starting on port :%s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
