package config

import "testing"

func TestLoadFromEnvironmentWithoutDotEnv(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("PORT", "10000")
	t.Setenv("CLERK_SECRET_KEY", "test-secret")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.example.com")

	cfg := Load()
	if cfg.DatabaseURL != "postgres://example" || cfg.Port != "10000" || cfg.ClerkSecretKey != "test-secret" || cfg.CORSAllowedOrigins != "https://app.example.com" {
		t.Fatalf("environment settings were not loaded: %#v", cfg)
	}
}
