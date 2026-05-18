package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ConfigInterface interface {
	GetPort() string
	GetDBUrl() string
	GetGoogleClientID() string
	GetJWTSecret() string
	GetRedisAddr() string
}

type Config struct {
	port           string
	dbUrl          string
	jwtSecret      string
	googleClientID string
	redisAddr      string
}

func NewConfig() *Config {
	godotenv.Load()

	return &Config{
		port:           os.Getenv("PORT"),
		dbUrl:          os.Getenv("DB_URL"),
		jwtSecret:      os.Getenv("JWT_SECRET"),
		googleClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		redisAddr:      os.Getenv("REDIS_ADDR"),
	}
}

func (c *Config) GetPort() string {
	return c.port
}

func (c *Config) GetDBUrl() string {
	return c.dbUrl
}

func (c *Config) GetGoogleClientID() string {
	return c.googleClientID
}

func (c *Config) GetJWTSecret() string {
	return c.jwtSecret
}

func (c *Config) GetRedisAddr() string {
	return c.redisAddr
}