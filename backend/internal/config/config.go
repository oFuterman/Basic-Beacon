package config

import "os"

type Config struct {
	DatabaseURL  string
	JWTSecret    string
	CORSOrigins  string
	SendGridKey  string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	Environment  string
	FrontendURL  string
	// Stripe configuration
	StripeSecretKey      string
	StripeWebhookSecret  string
	StripeIndiePriceID   string
	StripeTeamPriceID    string
	StripeAgencyPriceID  string
}

func Load() *Config {
	return &Config{
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/light_house?sslmode=disable"),
		JWTSecret:           getEnv("JWT_SECRET", "change-me-in-production"),
		CORSOrigins:         getEnv("CORS_ORIGINS", "*"),
		SendGridKey:         getEnv("SENDGRID_API_KEY", ""),
		SMTPHost:            getEnv("SMTP_HOST", ""),
		SMTPPort:            getEnv("SMTP_PORT", "587"),
		SMTPUser:            getEnv("SMTP_USER", ""),
		SMTPPassword:        getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:            getEnv("SMTP_FROM", "alerts@lighthouse.local"),
		Environment:         getEnv("ENVIRONMENT", "development"),
		FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:3000"),
		StripeSecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		StripeIndiePriceID:  getEnv("STRIPE_INDIE_PRICE_ID", ""),
		StripeTeamPriceID:   getEnv("STRIPE_TEAM_PRICE_ID", ""),
		StripeAgencyPriceID: getEnv("STRIPE_AGENCY_PRICE_ID", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
