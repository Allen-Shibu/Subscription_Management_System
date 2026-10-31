package config

import (
	"os"
)

type Config struct {
	DatabaseURL         string
	Port                string
	ClerkSecretKey      string
	ClerkPublishableKey string
	DevelopmentMode     string
}

func Load() *Config {
	//setting up config

	config := Config{
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		Port:                os.Getenv("PORT"),
		ClerkSecretKey:      os.Getenv("CLERK_SECRET_KEY"),
		ClerkPublishableKey: os.Getenv("CLERK_PUBLISHER_KEY"),
		DevelopmentMode:     os.Getenv("DEVELOPMENT_MODE"),
	}

	return &config
}
