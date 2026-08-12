package config

import (
	"os"
	"strconv"
	"strings"

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
	GetMinioEndpoint() string
	GetMinioPublicEndpoint() string
	GetMinioAccessKey() string
	GetMinioSecretKey() string
	GetMinioBucket() string
	GetClientOrigins() []string
}

type Config struct {
	port                       string
	dbUrl                      string
	jwtSecret                  string
	googleClientID             string
	redisAddr                  string
	isProduction               bool
	allowCredentials           bool
	minioEndpoint              string
	minioPublicEndpoint        string
	minioAccessKey             string
	minioSecretKey             string
	minioBucket                string
	clientOrigins              []string
}

func NewConfig() *Config {
	godotenv.Load()

	isProduction, _ := strconv.ParseBool(os.Getenv("IS_PRODUCTION"))
	allowCredentials, _ := strconv.ParseBool(os.Getenv("ALLOW_CREDENTIALS"))
	clientOrigins := parseOrigins(os.Getenv("CLIENT_ORIGINS"))

	return &Config{
		port:                       os.Getenv("PORT"),
		dbUrl:                      os.Getenv("DB_URL"),
		jwtSecret:                  os.Getenv("JWT_SECRET"),
		googleClientID:             os.Getenv("GOOGLE_CLIENT_ID"),
		redisAddr:                  os.Getenv("REDIS_ADDR"),
		isProduction:               isProduction,
		allowCredentials:           allowCredentials,
		minioEndpoint:              os.Getenv("MINIO_ENDPOINT"),
		minioPublicEndpoint:        os.Getenv("MINIO_PUBLIC_ENDPOINT"),
		minioAccessKey:             os.Getenv("MINIO_ACCESS_KEY"),
		minioSecretKey:             os.Getenv("MINIO_SECRET_KEY"),
		minioBucket:                os.Getenv("MINIO_BUCKET"),
		clientOrigins: clientOrigins,
	}
}

func (c *Config) GetPort() string                 { return c.port }
func (c *Config) GetDBUrl() string                { return c.dbUrl }
func (c *Config) GetGoogleClientID() string       { return c.googleClientID }
func (c *Config) GetJWTSecret() string            { return c.jwtSecret }
func (c *Config) GetRedisAddr() string            { return c.redisAddr }
func (c *Config) IsProduction() bool              { return c.isProduction }
func (c *Config) AllowCredentials() bool          { return c.allowCredentials }
func (c *Config) GetMinioEndpoint() string        { return c.minioEndpoint }
func (c *Config) GetMinioPublicEndpoint() string  { return c.minioPublicEndpoint }
func (c *Config) GetMinioAccessKey() string       { return c.minioAccessKey }
func (c *Config) GetMinioSecretKey() string       { return c.minioSecretKey }
func (c *Config) GetMinioBucket() string          { return c.minioBucket }
func (c *Config) GetClientOrigins() []string      { return c.clientOrigins }

func parseOrigins(raw string) []string {
	if raw == "" {
		return []string{"http://localhost:3000"}
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			origins = append(origins, p)
		}
	}
	return origins
}