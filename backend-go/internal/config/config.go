package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	AppEnv  string
	AppHost string
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret      string
	JWTExpireHours int

	AdminLogin    string
	AdminPassword string
	AdminFullName string

	UploadDir   string
	MaxUploadMB int

	TelegramBotToken string
	TelegramDebug    bool

	CORSOrigins []string
}

// Load reads configuration from .env (if present) and process environment.
func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppHost: getEnv("APP_HOST", "0.0.0.0"),
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "crm"),
		DBPassword: getEnv("DB_PASSWORD", "crm_password"),
		DBName:     getEnv("DB_NAME", "crm_real"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:      getEnv("JWT_SECRET", "change-me"),
		JWTExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 72),

		AdminLogin:    getEnv("ADMIN_LOGIN", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin"),
		AdminFullName: getEnv("ADMIN_FULL_NAME", "Administrator"),

		UploadDir:   getEnv("UPLOAD_DIR", "uploads"),
		MaxUploadMB: getEnvInt("MAX_UPLOAD_MB", 50),

		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramDebug:    getEnvBool("TELEGRAM_DEBUG", false),

		CORSOrigins: splitCSV(getEnv("CORS_ORIGINS", "*")),
	}
	return cfg
}

// PostgresDSN returns a libpq style DSN.
func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func (c *Config) IsProduction() bool {
	return strings.EqualFold(c.AppEnv, "production")
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}
