package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	RabbitMQ struct {
		URL string
	}
	JWT struct {
		Secret string
	}
	EncryptionKey string
	Alertmanager  struct {
		URL string
	}
	Email struct {
		Host     string
		Port     int
		From     string
		Password string
	}
}

func LoadConfig() *Config {
	// Load .env file if it exists
	godotenv.Load()

	cfg := &Config{}

	// Database Config
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvAsInt("DB_PORT", 5432)
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.DBName = getEnv("DB_NAME", "db_monitron")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// RabbitMQ Config
	cfg.RabbitMQ.URL = getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	// JWT Config
	cfg.JWT.Secret = getEnv("JWT_SECRET", "supersecretjwtkey")

	// Encryption Key
	cfg.EncryptionKey = getEnv("ENCRYPTION_KEY", "averysecretkey1234567890123456789012") // 32 bytes for AES-256

	// Alertmanager Config
	cfg.Alertmanager.URL = getEnv("ALERTMANAGER_URL", "http://localhost:9093")

	// Email Config
	cfg.Email.Host = getEnv("EMAIL_HOST", "smtp.mailtrap.io")
	cfg.Email.Port = getEnvAsInt("EMAIL_PORT", 2525)
	cfg.Email.From = getEnv("EMAIL_FROM", "no-reply@monitron.com")
	cfg.Email.Password = getEnv("EMAIL_PASSWORD", "your_email_password")

	return cfg
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
