package com.welcomeuniversity.provas.dto.subject;

import jakarta.validation.constraints.NotBlank;

public record CreateSubjectRequest(
    @NotBlank String name
) {
}
