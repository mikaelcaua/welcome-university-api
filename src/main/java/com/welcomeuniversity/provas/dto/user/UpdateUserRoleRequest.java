package com.welcomeuniversity.provas.dto.user;

import com.welcomeuniversity.provas.model.Role;

import jakarta.validation.constraints.NotNull;

public record UpdateUserRoleRequest(@NotNull Role role) {
}
