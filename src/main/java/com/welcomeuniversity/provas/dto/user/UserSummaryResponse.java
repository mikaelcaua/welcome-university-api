package com.welcomeuniversity.provas.dto.user;

import com.welcomeuniversity.provas.model.AppUser;
import com.welcomeuniversity.provas.model.Role;

public record UserSummaryResponse(
    Long id,
    String name,
    String email,
    Role role
) {

    public static UserSummaryResponse from(AppUser user) {
        if (user == null) {
            return null;
        }

        return new UserSummaryResponse(user.getId(), user.getName(), user.getEmail(), user.getRole());
    }
}
