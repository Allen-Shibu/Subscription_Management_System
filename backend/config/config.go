package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL         string
	Port                string
	ClerkSecretKey      string
	ClerkPublishableKey string
	DevelopmentMode     string
	SMTPHost            string
	SMTPPort            string
	SMTPUsername        string
	SMTPPassword        string
	SMTPFrom            string
	CORSAllowedOrigins  string
}

func Load() *Config {
	// Load .env file if present; ignore error if missing so environment variables can be used
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] No .env file found; using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587"
	}

	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:5173,http://127.0.0.1:5173"
	}

	config := Config{
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		Port:                port,
		ClerkSecretKey:      os.Getenv("CLERK_SECRET_KEY"),
		ClerkPublishableKey: os.Getenv("CLERK_PUBLISHABLE_KEY"),
		DevelopmentMode:     os.Getenv("DEVELOPMENT_MODE"),
		SMTPHost:            os.Getenv("SMTP_HOST"),
		SMTPPort:            smtpPort,
		SMTPUsername:        os.Getenv("SMTP_USERNAME"),
		SMTPPassword:        os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:            os.Getenv("SMTP_FROM"),
		CORSAllowedOrigins:  corsOrigins,
	}

	return &config
}
