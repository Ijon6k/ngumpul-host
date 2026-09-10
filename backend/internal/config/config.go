package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv         string
	AppURL         string
	Port           string
	DatabaseURL    string
	SessionSecret  string
	SessionSecure  bool
	StorageDriver  string
	StoragePath    string
	SMTPHost       string
	SMTPPort       int
	SMTPUsername   string
	SMTPPassword   string
	SMTPFrom       string
	BootstrapAdmin bool
	AdminEmail     string
	AdminUsername  string
	AdminPassword  string
	AdminName      string
}

func Load() *Config {
	sessionSecure := false
	if secureEnv := os.Getenv("SESSION_SECURE"); secureEnv == "true" || secureEnv == "1" {
		sessionSecure = true
	}

	bootstrapAdmin := true
	if bootstrapEnv := os.Getenv("BOOTSTRAP_ADMIN"); bootstrapEnv == "false" || bootstrapEnv == "0" {
		bootstrapAdmin = false
	}

	smtpPort := 587
	if portStr := os.Getenv("SMTP_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			smtpPort = p
		}
	}

	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://ngumpul:ngumpul_secret@localhost:5432/ngumpul_db?sslmode=disable")
	sessionSecret := getEnv("SESSION_SECRET", "ngumpul-host-default-insecure-secret-please-change-in-env-32chars")
	storagePath := getEnv("STORAGE_PATH", "./uploads")

	return &Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		AppURL:         getEnv("APP_URL", "http://localhost"),
		Port:           port,
		DatabaseURL:    dbURL,
		SessionSecret:  sessionSecret,
		SessionSecure:  sessionSecure,
		StorageDriver:  getEnv("STORAGE_DRIVER", "local"),
		StoragePath:    storagePath,
		SMTPHost:       os.Getenv("SMTP_HOST"),
		SMTPPort:       smtpPort,
		SMTPUsername:   os.Getenv("SMTP_USERNAME"),
		SMTPPassword:   os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:       getEnv("SMTP_FROM", "Ngumpul Host <no-reply@example.com>"),
		BootstrapAdmin: bootstrapAdmin,
		AdminEmail:     getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminUsername:  getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", "change_me_immediately_1234"),
		AdminName:      getEnv("ADMIN_NAME", "Site Administrator"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
