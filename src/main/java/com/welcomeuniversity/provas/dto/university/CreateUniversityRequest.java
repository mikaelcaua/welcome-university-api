package com.welcomeuniversity.provas.dto.university;

import jakarta.validation.constraints.NotBlank;

public record CreateUniversityRequest(
    @NotBlank String name
) {
}
