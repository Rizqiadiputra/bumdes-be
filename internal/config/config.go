package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	JWTSecret   string
	JWTExpHours int

	// BillingDueReminderDays adalah H- berapa hari sebelum tanggal jatuh
	// tempo tagihan dibuat oleh cron billing.
	BillingDueReminderDays int
}

func Load() *Config {
	_ = godotenv.Load()

	expHours, err := strconv.Atoi(getEnv("JWT_EXP_HOURS", "24"))
	if err != nil {
		expHours = 24
	}

	billingDueReminderDays, err := strconv.Atoi(getEnv("BILLING_DUE_REMINDER_DAYS", "7"))
	if err != nil {
		billingDueReminderDays = 7
	}

	return &Config{
		AppPort:                getEnv("APP_PORT", "8080"),
		DBHost:                 getEnv("DB_HOST", "localhost"),
		DBPort:                 getEnv("DB_PORT", "5432"),
		DBUser:                 getEnv("DB_USER", "postgres"),
		DBPassword:             getEnv("DB_PASSWORD", "postgres"),
		DBName:                 getEnv("DB_NAME", "bumdes"),
		DBSSLMode:              getEnv("DB_SSLMODE", "disable"),
		JWTSecret:              getEnv("JWT_SECRET", "change-me-secret"),
		JWTExpHours:            expHours,
		BillingDueReminderDays: billingDueReminderDays,
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
