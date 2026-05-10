package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"github.com/mikaelcaua/welcome-university-api/internal/config"
	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
)

type UploadPayload struct {
	OriginalFilename string
	ContentType      string
	Bytes            []byte
}

type StoredObject struct {
	Key string
	URL string
}

type ObjectStorage interface {
	UploadExam(ctx context.Context, payload UploadPayload, subjectID int64) (StoredObject, error)
	DeleteObjectByKey(ctx context.Context, storageKey *string) error
}

type S3StorageService struct {
	client         *s3.Client
	bucketName     string
	endpoint       string
	publicEndpoint string
}

func NewS3StorageService(ctx context.Context, appConfig config.Config) (*S3StorageService, error) {
	awsConfig, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(appConfig.S3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(appConfig.S3AccessKey, appConfig.S3SecretKey, "")),
	)
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(awsConfig, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(appConfig.S3Endpoint)
		options.UsePathStyle = true
	})
	service := &S3StorageService{
		client:         client,
		bucketName:     appConfig.S3BucketName,
		endpoint:       appConfig.S3Endpoint,
		publicEndpoint: appConfig.S3PublicEndpoint,
	}
	if err := service.ensureBucketExists(ctx); err != nil {
		return nil, err
	}
	return service, nil
}

func (service *S3StorageService) UploadExam(ctx context.Context, payload UploadPayload, subjectID int64) (StoredObject, error) {
	safeFilename := sanitizeFilename(payload.OriginalFilename)
	objectKey := fmt.Sprintf("subjects/%d/%s-%s", subjectID, uuid.NewString(), safeFilename)
	contentType := resolveContentType(payload.ContentType, safeFilename)
	_, err := service.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(service.bucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
		Body:        bytes.NewReader(payload.Bytes),
	})
	if err != nil {
		return StoredObject{}, httpx.NewHTTPError(http.StatusBadGateway, "Falha ao enviar arquivo para o S3.")
	}
	return StoredObject{Key: objectKey, URL: service.buildPublicURL(objectKey)}, nil
}

func (service *S3StorageService) DeleteObjectByKey(ctx context.Context, storageKey *string) error {
	if storageKey == nil || strings.TrimSpace(*storageKey) == "" {
		return nil
	}
	_, err := service.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(service.bucketName),
		Key:    aws.String(*storageKey),
	})
	if err != nil {
		return httpx.NewHTTPError(http.StatusBadGateway, "Falha ao remover arquivo do S3.")
	}
	return nil
}

func (service *S3StorageService) ensureBucketExists(ctx context.Context) error {
	_, err := service.client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(service.bucketName)})
	if err == nil {
		return nil
	}
	_, err = service.client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(service.bucketName)})
	return err
}

func (service *S3StorageService) buildPublicURL(objectKey string) string {
	baseURL := strings.TrimRight(service.publicEndpoint, "/")
	if baseURL == "" {
		baseURL = strings.TrimRight(service.endpoint, "/")
	}
	if baseURL == "" {
		return fmt.Sprintf("s3://%s/%s", service.bucketName, objectKey)
	}
	return fmt.Sprintf("%s/%s/%s", baseURL, service.bucketName, objectKey)
}

type UploadOptimizer struct{}

func NewUploadOptimizer() *UploadOptimizer {
	return &UploadOptimizer{}
}

func (optimizer *UploadOptimizer) Optimize(payload UploadPayload) UploadPayload {
	lowerFilename := strings.ToLower(payload.OriginalFilename)
	lowerContentType := strings.ToLower(payload.ContentType)
	if lowerContentType == "image/jpeg" || lowerContentType == "image/jpg" || strings.HasSuffix(lowerFilename, ".jpg") || strings.HasSuffix(lowerFilename, ".jpeg") {
		return optimizer.optimizeJPEG(payload)
	}
	if lowerContentType == "image/png" || strings.HasSuffix(lowerFilename, ".png") {
		return optimizer.optimizePNG(payload)
	}
	return payload
}

func (optimizer *UploadOptimizer) optimizeJPEG(payload UploadPayload) UploadPayload {
	sourceImage, _, err := image.Decode(bytes.NewReader(payload.Bytes))
	if err != nil {
		return payload
	}
	var output bytes.Buffer
	if err := jpeg.Encode(&output, sourceImage, &jpeg.Options{Quality: 78}); err != nil {
		return payload
	}
	if output.Len() >= len(payload.Bytes) {
		return payload
	}
	payload.Bytes = output.Bytes()
	payload.ContentType = "image/jpeg"
	return payload
}

func (optimizer *UploadOptimizer) optimizePNG(payload UploadPayload) UploadPayload {
	sourceImage, _, err := image.Decode(bytes.NewReader(payload.Bytes))
	if err != nil {
		return payload
	}
	var output bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	if err := encoder.Encode(&output, sourceImage); err != nil {
		return payload
	}
	if output.Len() >= len(payload.Bytes) {
		return payload
	}
	payload.Bytes = output.Bytes()
	payload.ContentType = "image/png"
	return payload
}

func NewUploadPayload(originalFilename string, contentType string, reader io.Reader) (UploadPayload, error) {
	fileBytes, err := io.ReadAll(reader)
	if err != nil {
		return UploadPayload{}, httpx.NewHTTPError(http.StatusBadRequest, "Falha ao ler arquivo enviado.")
	}
	return UploadPayload{
		OriginalFilename: originalFilename,
		ContentType:      contentType,
		Bytes:            fileBytes,
	}, nil
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
