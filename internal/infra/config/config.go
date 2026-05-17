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
	MinioEndpoint string
	MinioPublicEndpoint string
	MinioBucketName string
	MinioAccessKey string
	MinioSecretKey string
	MinioRegion string
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
		MinioEndpoint:                getEnvAny([]string{"MINIO_ENDPOINT", "STORAGE_ENDPOINT"}, "http://minio:9000"),
		MinioPublicEndpoint:          getEnvAny([]string{"MINIO_PUBLIC_ENDPOINT", "S3_PUBLIC_ENDPOINT", "STORAGE_PUBLIC_ENDPOINT"}, "http://localhost:9002"),
		MinioBucketName:              getEnvAny([]string{"MINIO_BUCKET_NAME", "S3_BUCKET_NAME", "STORAGE_BUCKET_NAME"}, "exams-bucket"),
		MinioAccessKey:               getEnvAny([]string{"MINIO_ROOT_USER", "STORAGE_ACCESS_KEY"}, ""),
		MinioSecretKey:               getEnvAny([]string{"MINIO_ROOT_PASSWORD", "STORAGE_SECRET_KEY"}, ""),
		MinioRegion:                  getEnvAny([]string{"MINIO_REGION", "AWS_REGION", "STORAGE_REGION"}, "us-east-1"),
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
	if config.MinioAccessKey == "" || config.MinioSecretKey == "" || config.MinioBucketName == "" {
		return fmt.Errorf("MINIO_ROOT_USER, MINIO_ROOT_PASSWORD and MINIO_BUCKET_NAME are required")
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

func getEnvAny(names []string, fallback string) string {
	for _, name := range names {
		if value := os.Getenv(name); value != "" {
			return value
		}
	}
	return fallback
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
