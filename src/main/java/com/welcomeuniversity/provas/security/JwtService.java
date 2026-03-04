package com.welcomeuniversity.provas.security;

import java.nio.charset.StandardCharsets;
import java.util.Date;

import javax.crypto.SecretKey;

import org.springframework.stereotype.Service;

import com.welcomeuniversity.provas.config.JwtProperties;
import com.welcomeuniversity.provas.model.AppUser;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.Keys;

@Service
public class JwtService {

    private static final String ACCESS_TOKEN_TYPE = "access";
    private static final String REFRESH_TOKEN_TYPE = "refresh";

    private final JwtProperties properties;
    private final SecretKey signingKey;

    public JwtService(JwtProperties properties) {
        this.properties = properties;
        this.signingKey = Keys.hmacShaKeyFor(properties.jwtSecret().getBytes(StandardCharsets.UTF_8));
    }

    public String generateAccessToken(AppUser user) {
        return generateToken(user, ACCESS_TOKEN_TYPE, properties.accessTokenExpiration());
    }

    public String generateRefreshToken(AppUser user) {
        return generateToken(user, REFRESH_TOKEN_TYPE, properties.refreshTokenExpiration());
    }

    public String extractEmailFromAccessToken(String token) {
        return extractEmail(token, ACCESS_TOKEN_TYPE);
    }

    public String extractEmailFromRefreshToken(String token) {
        return extractEmail(token, REFRESH_TOKEN_TYPE);
    }

    public long getAccessTokenExpiration() {
        return properties.accessTokenExpiration();
    }

    private String generateToken(AppUser user, String tokenType, long expirationSeconds) {
        Date issuedAt = new Date();
        Date expiration = new Date(issuedAt.getTime() + expirationSeconds * 1000);

        return Jwts.builder()
            .subject(user.getEmail())
            .claim("uid", user.getId())
            .claim("role", user.getRole().name())
            .claim("type", tokenType)
            .issuedAt(issuedAt)
            .expiration(expiration)
            .signWith(signingKey)
            .compact();
    }

    private String extractEmail(String token, String expectedType) {
        Claims claims = Jwts.parser()
            .verifyWith(signingKey)
            .build()
            .parseSignedClaims(token)
            .getPayload();

        String tokenType = claims.get("type", String.class);
        if (!expectedType.equals(tokenType)) {
            throw new IllegalArgumentException("Tipo de token invalido.");
        }

        return claims.getSubject();
    }
}
