package main

import (
	"log"
	"net/http"
	"strconv"

	"recur_backend/config"
	"recur_backend/internals/handlers"
	"recur_backend/internals/middleware"
	"recur_backend/internals/repository"
	invoicesvc "recur_backend/internals/services/invoice"
	"recur_backend/pkg/database"
	"recur_backend/pkg/email"
	"recur_backend/pkg/types"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupRouter(pool *pgxpool.Pool, mailer email.Sender, allowedOrigins string) *gin.Engine {
	r := gin.Default()

	// Apply CORS middleware
	r.Use(middleware.CORSMiddleware(allowedOrigins))

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
				if rc, ok := claims.Custom.(*middleware.RoleClaims); ok {
					role = rc.Role
				}
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

		staff := api.Group("/staff")
		staff.Use(middleware.RequireAnyRole(types.RoleAdmin, types.RoleStaff))

		productHandler := handlers.NewProductHandler(repository.NewProductRepository(pool))
		productHandler.RegisterRoutes(api, admin)

		subscriptionRepo := repository.NewSubscriptionRepository(pool)
		invoiceService := invoicesvc.New(subscriptionRepo, mailer)
		subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionRepo, invoiceService)
		subscriptionHandler.RegisterRoutes(staff, admin)

		taxRuleHandler := handlers.NewTaxRuleHandler(repository.NewTaxRuleRepository(pool))
		taxRuleHandler.RegisterRoutes(api, admin)

		handlers.RegisterStaffRoutes(admin)
	}

	return r
}

func main() {
	// Load required .env file.
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

	smtpPort, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil || smtpPort == 0 {
		smtpPort = 587
	}
	mailer := email.NewSMTP(email.Config{
		Host:     cfg.SMTPHost,
		Port:     smtpPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
	})
	if cfg.SMTPHost == "" {
		log.Println("[WARN] SMTP_HOST is not set; invoice emails will fail to send")
	}

	r := setupRouter(pool, mailer, cfg.CORSAllowedOrigins)

	log.Printf("[INFO] Server starting on port :%s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
