package config

import "os"

type Config struct {
	Port           string
	DBUrl          string
	JWTSecret      string
	GoogleClientID string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DBUrl:          os.Getenv("DB_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		GoogleClientID: os.Getenv("GOOGLE_CLIENT_ID"),
	}
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}