package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
	DatabaseURL string
	DatabaseUser string
	DatabasePassword string
	DatabaseName string

	JWTSecret string
	AccessTokenExpirationSeconds int64
	RefreshTokenExpirationSeconds int64
	MaxUploadSizeMB int64
	S3Endpoint string
	S3PublicEndpoint string
	S3BucketName string
	S3AccessKey string
	S3SecretKey string
	S3Region string
}

func Load() Config {
	return Config{
		ServerPort:                    getEnv("API_PORT", "8080"),
		DatabaseURL:                   getEnv("DATABASE_URL", ""),
		DatabaseUser:                  getEnv("DB_USER", ""),
		DatabasePassword:              getEnv("DB_PASSWORD", ""),
		DatabaseName:                  getEnv("DB_NAME", ""),
		JWTSecret:                     getEnv("API_JWT_SECRET", ""),
		AccessTokenExpirationSeconds:  getEnvInt64("API_ACCESS_TOKEN_EXPIRATION_SECONDS", 900),
		RefreshTokenExpirationSeconds: getEnvInt64("API_REFRESH_TOKEN_EXPIRATION_SECONDS", 604800),
		MaxUploadSizeMB:               getEnvInt64("API_MAX_UPLOAD_MB", 25),
		S3Endpoint:                    getEnv("STORAGE_ENDPOINT", "http://minio:9000"),
		S3PublicEndpoint:              getEnv("STORAGE_PUBLIC_ENDPOINT", "http://localhost:9002"),
		S3BucketName:                  getEnv("STORAGE_BUCKET_NAME", "exams-bucket"),
		S3AccessKey:                   getEnv("STORAGE_ACCESS_KEY", ""),
		S3SecretKey:                   getEnv("STORAGE_SECRET_KEY", ""),
		S3Region:                      getEnv("STORAGE_REGION", "us-east-1"),
	}
}

func (config Config) Validate() error {
	if config.JWTSecret == "" {
		return fmt.Errorf("API_JWT_SECRET is required")
	}
	if len(config.JWTSecret) < 32 {
		return fmt.Errorf("API_JWT_SECRET must have at least 32 characters")
	}
	if config.DatabaseUser == "" || config.DatabasePassword == "" || config.DatabaseName == "" {
		return fmt.Errorf("DB_NAME, DB_USER and DB_PASSWORD are required")
	}
	if config.S3AccessKey == "" || config.S3SecretKey == "" || config.S3BucketName == "" {
		return fmt.Errorf("STORAGE_ACCESS_KEY, STORAGE_SECRET_KEY and STORAGE_BUCKET_NAME are required")
	}
	return nil
}

func getEnv(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt64(name string, fallback int64) int64 {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsedValue
}
