package com.welcomeuniversity.provas.service;

import java.util.Locale;
import java.util.UUID;

import jakarta.annotation.PostConstruct;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import com.welcomeuniversity.provas.config.S3Properties;

import software.amazon.awssdk.core.sync.RequestBody;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.CreateBucketRequest;
import software.amazon.awssdk.services.s3.model.HeadBucketRequest;
import software.amazon.awssdk.services.s3.model.NoSuchBucketException;
import software.amazon.awssdk.services.s3.model.PutObjectRequest;
import software.amazon.awssdk.services.s3.model.S3Exception;

@Service
public class S3Service {

    private final S3Client s3Client;
    private final S3Properties properties;

    public S3Service(S3Client s3Client, S3Properties properties) {
        this.s3Client = s3Client;
        this.properties = properties;
    }

    @PostConstruct
    void ensureBucketExists() {
        try {
            s3Client.headBucket(HeadBucketRequest.builder().bucket(properties.bucketName()).build());
        } catch (NoSuchBucketException ex) {
            s3Client.createBucket(CreateBucketRequest.builder().bucket(properties.bucketName()).build());
        } catch (S3Exception ex) {
            if (ex.statusCode() == 404) {
                s3Client.createBucket(CreateBucketRequest.builder().bucket(properties.bucketName()).build());
                return;
            }
            throw ex;
        }
    }

    public StoredObject uploadExam(UploadPayload payload, Long subjectId) {
        String originalFilename = payload.originalFilename() == null ? "exam" : payload.originalFilename();
        String sanitizedFilename = originalFilename.replaceAll("[^a-zA-Z0-9._-]", "_");
        String objectKey = "subjects/%d/%s-%s".formatted(subjectId, UUID.randomUUID(), sanitizedFilename);
        String contentType = resolveContentType(payload.contentType(), sanitizedFilename);

        try {
            s3Client.putObject(
                PutObjectRequest.builder()
                    .bucket(properties.bucketName())
                    .key(objectKey)
                    .contentType(contentType)
                    .build(),
                RequestBody.fromBytes(payload.bytes())
            );
        } catch (S3Exception ex) {
            throw new ResponseStatusException(HttpStatus.BAD_GATEWAY, "Falha ao enviar arquivo para o S3.");
        }

        return new StoredObject(objectKey, buildObjectUrl(objectKey));
    }

    public String buildPublicUrlFromStorageKey(String storageKey) {
        if (storageKey == null || storageKey.isBlank()) {
            return null;
        }
        return buildObjectUrl(storageKey);
    }

    private String buildObjectUrl(String objectKey) {
        String publicEndpoint = properties.publicEndpoint();
        if (publicEndpoint == null || publicEndpoint.isBlank()) {
            publicEndpoint = properties.endpoint();
        }

        if (publicEndpoint == null || publicEndpoint.isBlank()) {
            return "s3://%s/%s".formatted(properties.bucketName(), objectKey);
        }

        String baseEndpoint = publicEndpoint.endsWith("/")
            ? publicEndpoint.substring(0, publicEndpoint.length() - 1)
            : publicEndpoint;
        return "%s/%s/%s".formatted(baseEndpoint, properties.bucketName(), objectKey);
    }

    private String resolveContentType(String providedContentType, String filename) {
        String contentType = providedContentType;
        if (contentType != null && !contentType.isBlank()) {
            return contentType;
        }

        String lowerFilename = filename.toLowerCase(Locale.ROOT);
        if (lowerFilename.endsWith(".pdf")) {
            return "application/pdf";
        }
        if (lowerFilename.endsWith(".png")) {
            return "image/png";
        }
        if (lowerFilename.endsWith(".jpg") || lowerFilename.endsWith(".jpeg")) {
            return "image/jpeg";
        }
        if (lowerFilename.endsWith(".webp")) {
            return "image/webp";
        }
        if (lowerFilename.endsWith(".gif")) {
            return "image/gif";
        }
        return "application/octet-stream";
    }

    public record StoredObject(String key, String url) {
    }

    public record UploadPayload(String originalFilename, String contentType, byte[] bytes) {
    }
}
