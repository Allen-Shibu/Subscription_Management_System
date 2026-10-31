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
<<<<<<< HEAD
	SMTPHost            string
	SMTPPort            string
	SMTPUsername        string
	SMTPPassword        string
	SMTPFrom            string
=======
	CORSAllowedOrigins  string
>>>>>>> 2ff69f2 (db:connection with neon)
}

func Load() *Config {
	//setting up config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := Config{
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		Port:                os.Getenv("PORT"),
		ClerkSecretKey:      os.Getenv("CLERK_SECRET_KEY"),
		ClerkPublishableKey: os.Getenv("CLERK_PUBLISHABLE_KEY"),
		DevelopmentMode:     os.Getenv("DEVELOPMENT_MODE"),
<<<<<<< HEAD
		SMTPHost:            os.Getenv("SMTP_HOST"),
		SMTPPort:            os.Getenv("SMTP_PORT"),
		SMTPUsername:        os.Getenv("SMTP_USERNAME"),
		SMTPPassword:        os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:            os.Getenv("SMTP_FROM"),
=======
		CORSAllowedOrigins:  os.Getenv("CORS_ALLOWED_ORIGINS"),
>>>>>>> 2ff69f2 (db:connection with neon)
	}

	return &config
}
