package com.welcomeuniversity.provas.dto.user;

import java.time.Instant;

import com.welcomeuniversity.provas.model.AppUser;
import com.welcomeuniversity.provas.model.Role;

public record UserResponse(
    Long id,
    String name,
    String email,
    Role role,
    Instant createdAt
) {

    public static UserResponse from(AppUser user) {
        return new UserResponse(
            user.getId(),
            user.getName(),
            user.getEmail(),
            user.getRole(),
            user.getCreatedAt()
        );
    }
}
