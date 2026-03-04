package com.welcomeuniversity.provas.dto.course;

import jakarta.validation.constraints.NotBlank;

public record CreateCourseRequest(
    @NotBlank String name
) {
}
