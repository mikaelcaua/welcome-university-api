package com.welcomeuniversity.provas.dto.auth;

import com.welcomeuniversity.provas.dto.user.UserResponse;

public record AuthResponse(
    String accessToken,
    String refreshToken,
    String tokenType,
    long expiresInSeconds,
    UserResponse user
) {
}
