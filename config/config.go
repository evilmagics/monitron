package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	Port      string
	JWTSecret string
}

func (c Config) GenerateDatabaseURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

func Load() *Config {
	return &Config{
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "3306"),
		DBUser:    getEnv("DB_USER", "root"),
		DBPass:    getEnv("DB_PASS", ""),
		DBName:    getEnv("DB_NAME", "service_monitoring"),
		Port:      getEnv("APP_PORT", "8080"),
		JWTSecret: getEnv("APP_SECRET", "your-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
