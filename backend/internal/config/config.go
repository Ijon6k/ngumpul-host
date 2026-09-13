package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv        string
	AppURL        string
	Port          string
	DatabaseURL   string
	SessionSecret string
	SessionSecure bool
	StorageDriver string
	StoragePath   string
	S3Endpoint    string
	S3Region      string
	S3Bucket      string
	S3AccessKey   string
	S3SecretKey   string
	S3PublicURL   string
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	SMTPFrom      string
}

func Load() *Config {
	sessionSecure := false
	if secureEnv := os.Getenv("SESSION_SECURE"); secureEnv == "true" || secureEnv == "1" {
		sessionSecure = true
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
		AppEnv:        getEnv("APP_ENV", "development"),
		AppURL:        getEnv("APP_URL", "http://localhost"),
		Port:          port,
		DatabaseURL:   dbURL,
		SessionSecret: sessionSecret,
		SessionSecure: sessionSecure,
		StorageDriver: getEnv("STORAGE_DRIVER", "local"),
		StoragePath:   storagePath,
		S3Endpoint:    getEnv("S3_ENDPOINT", "http://seaweedfs:8333"),
		S3Region:      getEnv("S3_REGION", "us-east-1"),
		S3Bucket:      getEnv("S3_BUCKET", "ngumpul-uploads"),
		S3AccessKey:   os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:   os.Getenv("S3_SECRET_KEY"),
		S3PublicURL:   getEnv("S3_PUBLIC_URL", "/uploads"),
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      smtpPort,
		SMTPUsername:  os.Getenv("SMTP_USERNAME"),
		SMTPPassword:  os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:      getEnv("SMTP_FROM", "Ngumpul Host <no-reply@example.com>"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
