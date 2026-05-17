package storagerepository

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/storage"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/config"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/httpx"
)

type MinioStorageRepositoryImpl struct {
	client         *s3.Client
	bucketName     string
	endpoint       string
	publicEndpoint string
}

func NewMinioStorageRepositoryImpl(ctx context.Context, appConfig config.Config) (*MinioStorageRepositoryImpl, error) {
	awsConfig, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(appConfig.MinioRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(appConfig.MinioAccessKey, appConfig.MinioSecretKey, "")),
	)
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(awsConfig, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(appConfig.MinioEndpoint)
		options.UsePathStyle = true
	})
	repository := &MinioStorageRepositoryImpl{client: client, bucketName: appConfig.MinioBucketName, endpoint: appConfig.MinioEndpoint, publicEndpoint: appConfig.MinioPublicEndpoint}
	if err := repository.ensureBucketExists(ctx); err != nil {
		return nil, err
	}
	return repository, nil
}

func (repository *MinioStorageRepositoryImpl) UploadExam(ctx context.Context, payload storagecontract.UploadPayload, subjectID int64) (storagecontract.StoredObject, error) {
	safeFilename := sanitizeFilename(payload.OriginalFilename)
	objectKey := fmt.Sprintf("subjects/%d/%s-%s", subjectID, uuid.NewString(), safeFilename)
	contentType := resolveContentType(payload.ContentType, safeFilename)
	_, err := repository.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(repository.bucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
		Body:        bytes.NewReader(payload.Bytes),
	})
	if err != nil {
		return storagecontract.StoredObject{}, httpx.NewHTTPError(http.StatusBadGateway, "Falha ao enviar arquivo para o MinIO.")
	}
	return storagecontract.StoredObject{Key: objectKey, URL: repository.buildPublicURL(objectKey)}, nil
}

func (repository *MinioStorageRepositoryImpl) DeleteObjectByKey(ctx context.Context, storageKey *string) error {
	if storageKey == nil || strings.TrimSpace(*storageKey) == "" {
		return nil
	}
	_, err := repository.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(repository.bucketName), Key: aws.String(*storageKey)})
	if err != nil {
		return httpx.NewHTTPError(http.StatusBadGateway, "Falha ao remover arquivo do MinIO.")
	}
	return nil
}

func (repository *MinioStorageRepositoryImpl) ensureBucketExists(ctx context.Context) error {
	_, err := repository.client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(repository.bucketName)})
	if err == nil {
		return nil
	}
	_, err = repository.client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(repository.bucketName)})
	return err
}

func (repository *MinioStorageRepositoryImpl) buildPublicURL(objectKey string) string {
	baseURL := strings.TrimRight(repository.publicEndpoint, "/")
	if baseURL == "" {
		baseURL = strings.TrimRight(repository.endpoint, "/")
	}
	if baseURL == "" {
		return fmt.Sprintf("minio://%s/%s", repository.bucketName, objectKey)
	}
	return fmt.Sprintf("%s/%s/%s", baseURL, repository.bucketName, objectKey)
}

func sanitizeFilename(filename string) string {
	if strings.TrimSpace(filename) == "" {
		filename = "exam"
	}
	filename = filepath.Base(filename)
	return regexp.MustCompile(`[^a-zA-Z0-9._-]`).ReplaceAllString(filename, "_")
}

func resolveContentType(providedContentType string, filename string) string {
	if strings.TrimSpace(providedContentType) != "" {
		return providedContentType
	}
	lowerFilename := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lowerFilename, ".pdf"):
		return "application/pdf"
	case strings.HasSuffix(lowerFilename, ".png"):
		return "image/png"
	case strings.HasSuffix(lowerFilename, ".jpg"), strings.HasSuffix(lowerFilename, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(lowerFilename, ".webp"):
		return "image/webp"
	case strings.HasSuffix(lowerFilename, ".gif"):
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
