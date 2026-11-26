package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config menyimpan seluruh konfigurasi aplikasi.
type Config struct {
	AppPort       string
	EnableSwagger bool
	Database      DatabaseConfig
	JWT           JWTConfig
	CORS          CORSConfig
	RateLimit     RateLimitConfig
	Storage       StorageConfig
	MFA           MFAConfig
}

type MFAConfig struct {
	Required        bool
	Issuer          string
	EncryptionKey   string
	RateLimitVerify int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

type JWTConfig struct {
	Secret             string
	AccessExpireHours  int
	RefreshExpireHours int
}

type CORSConfig struct {
	AllowedOrigins string
}

type RateLimitConfig struct {
	API           int
	Login         int
	WindowMinutes int
}

type StorageConfig struct {
	Path        string
	MaxUploadMB int
}

// DSN mengembalikan MySQL connection string.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.Name,
	)
}

// Load membaca konfigurasi dari .env dan environment variables.
// Wajib dipanggil di awal aplikasi (main.go).
func Load() (*Config, error) {
	// Load .env (abaikan error jika file tidak ada)
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:       getEnv("PORT", getEnv("APP_PORT", "8000")),
		EnableSwagger: getEnvBool("ENABLE_SWAGGER", true),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", getEnv("MYSQLHOST", "127.0.0.1")),
			Port:     getEnv("DB_PORT", getEnv("MYSQLPORT", "3306")),
			Name:     getEnv("DB_DATABASE", getEnv("MYSQLDATABASE", "siarsip")),
			Username: getEnv("DB_USERNAME", getEnv("MYSQLUSER", "root")),
			Password: getEnv("DB_PASSWORD", getEnv("MYSQLPASSWORD", "")),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", ""),
			AccessExpireHours:  getEnvInt("JWT_ACCESS_EXPIRATION_HOURS", 7),
			RefreshExpireHours: getEnvInt("JWT_REFRESH_EXPIRATION_HOURS", 30),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		},
		RateLimit: RateLimitConfig{
			API:           getEnvInt("RATE_LIMIT_API", 60),
			Login:         getEnvInt("RATE_LIMIT_LOGIN", 5),
			WindowMinutes: getEnvInt("RATE_LIMIT_WINDOW_MINUTES", 1),
		},
		Storage: StorageConfig{
			Path:        getEnv("STORAGE_PATH", "./storage/dokumen"),
			MaxUploadMB: getEnvInt("MAX_UPLOAD_MB", 10),
		},
		MFA: MFAConfig{
			Required:        getEnvBool("MFA_REQUIRED", true),
			Issuer:          getEnv("MFA_ISSUER", "SiARSIP"),
			EncryptionKey:   getEnv("MFA_ENCRYPTION_KEY", "siarsip-secret-key-32bytes-mfa!"),
			RateLimitVerify: getEnvInt("RATE_LIMIT_MFA", 5),
		},
	}

	// Validasi field required
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return fallback
}
