package com.welcomeuniversity.provas.config;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "app.s3")
public record S3Properties(
    String endpoint,
    String publicEndpoint,
    String bucketName,
    String accessKey,
    String secretKey,
    String region
) {
}
