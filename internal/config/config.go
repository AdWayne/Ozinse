package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	RefreshSecret string
	BaseURL       string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:          getEnv("PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "ozinse"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		RefreshSecret: getEnv("REFRESH_SECRET", ""),
		BaseURL:       getEnv("BASE_URL", "http://localhost:8080"),
	}
}

func (c *Config) DatabaseURL() string {
	return "postgres://" + c.DBUser + ":" + c.DBPassword +
		"@" + c.DBHost + ":" + c.DBPort +
		"/" + c.DBName + "?sslmode=disable"
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
