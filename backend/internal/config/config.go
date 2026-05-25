package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigInterface interface {
	GetPort() string
	GetDBUrl() string
	GetGoogleClientID() string
	GetJWTSecret() string
	GetRedisAddr() string
	IsProduction() bool
	AllowCredentials() bool
}

type Config struct {
	port             string
	dbUrl            string
	jwtSecret        string
	googleClientID   string
	redisAddr        string
	isProduction     bool
	allowCredentials bool
}

func NewConfig() *Config {
	godotenv.Load()

	isProduction, _ := strconv.ParseBool(
		os.Getenv("IS_PRODUCTION"),
	)

	allowCredentials, _ := strconv.ParseBool(
		os.Getenv("ALLOW_CREDENTIALS"),
	)

	return &Config{
		port:             os.Getenv("PORT"),
		dbUrl:            os.Getenv("DB_URL"),
		jwtSecret:        os.Getenv("JWT_SECRET"),
		googleClientID:   os.Getenv("GOOGLE_CLIENT_ID"),
		redisAddr:        os.Getenv("REDIS_ADDR"),
		isProduction:     isProduction,
		allowCredentials: allowCredentials,
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

func (c *Config) IsProduction() bool {
	return c.isProduction
}

func (c *Config) AllowCredentials() bool {
	return c.allowCredentials
}
