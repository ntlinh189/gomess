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
	GetFacebookAppID() string
	GetFacebookAppSecret() string
	GetJWTSecret() string
	GetRedisAddr() string
	IsProduction() bool
	AllowCredentials() bool
	GetMinioEndpoint() string
	GetMinioAccessKey() string
	GetMinioSecretKey() string
	GetMinioBucket() string
	MinioUseSSL() bool
	WorkerCleanupIntervalHours() int
}

type Config struct {
	port                       string
	dbUrl                      string
	jwtSecret                  string
	googleClientID             string
	facebookAppID              string
	facebookAppSecret          string
	redisAddr                  string
	isProduction               bool
	allowCredentials           bool
	minioEndpoint              string
	minioAccessKey             string
	minioSecretKey             string
	minioBucket                string
	minioUseSSL                bool
	workerCleanupIntervalHours int
}

func NewConfig() *Config {
	godotenv.Load()

	isProduction, _ := strconv.ParseBool(os.Getenv("IS_PRODUCTION"))
	allowCredentials, _ := strconv.ParseBool(os.Getenv("ALLOW_CREDENTIALS"))
	minioUseSSL, _ := strconv.ParseBool(os.Getenv("MINIO_USE_SSL"))
	workerCleanupIntervalHours, err := strconv.Atoi(os.Getenv("WORKER_CLEANUP_INTERVAL_HOURS"))
	if err != nil || workerCleanupIntervalHours <= 0 {
		workerCleanupIntervalHours = 6
	}

	return &Config{
		port:                       os.Getenv("PORT"),
		dbUrl:                      os.Getenv("DB_URL"),
		jwtSecret:                  os.Getenv("JWT_SECRET"),
		googleClientID:             os.Getenv("GOOGLE_CLIENT_ID"),
		facebookAppID:              os.Getenv("FACEBOOK_APP_ID"),
		facebookAppSecret:          os.Getenv("FACEBOOK_APP_SECRET"),
		redisAddr:                  os.Getenv("REDIS_ADDR"),
		isProduction:               isProduction,
		allowCredentials:           allowCredentials,
		minioEndpoint:              os.Getenv("MINIO_ENDPOINT"),
		minioAccessKey:             os.Getenv("MINIO_ACCESS_KEY"),
		minioSecretKey:             os.Getenv("MINIO_SECRET_KEY"),
		minioBucket:                os.Getenv("MINIO_BUCKET"),
		minioUseSSL:                minioUseSSL,
		workerCleanupIntervalHours: workerCleanupIntervalHours,
	}
}

func (c *Config) GetPort() string                 { return c.port }
func (c *Config) GetDBUrl() string                { return c.dbUrl }
func (c *Config) GetGoogleClientID() string       { return c.googleClientID }
func (c *Config) GetFacebookAppID() string        { return c.facebookAppID }
func (c *Config) GetFacebookAppSecret() string    { return c.facebookAppSecret }
func (c *Config) GetJWTSecret() string            { return c.jwtSecret }
func (c *Config) GetRedisAddr() string            { return c.redisAddr }
func (c *Config) IsProduction() bool              { return c.isProduction }
func (c *Config) AllowCredentials() bool          { return c.allowCredentials }
func (c *Config) GetMinioEndpoint() string        { return c.minioEndpoint }
func (c *Config) GetMinioAccessKey() string       { return c.minioAccessKey }
func (c *Config) GetMinioSecretKey() string       { return c.minioSecretKey }
func (c *Config) GetMinioBucket() string          { return c.minioBucket }
func (c *Config) MinioUseSSL() bool               { return c.minioUseSSL }
func (c *Config) WorkerCleanupIntervalHours() int { return c.workerCleanupIntervalHours }
