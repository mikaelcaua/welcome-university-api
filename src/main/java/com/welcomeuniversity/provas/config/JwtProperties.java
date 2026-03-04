package com.welcomeuniversity.provas.config;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "app.security")
public record JwtProperties(
    String jwtSecret,
    long accessTokenExpiration,
    long refreshTokenExpiration
) {
}
