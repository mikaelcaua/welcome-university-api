package com.welcomeuniversity.provas.dto.state;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

public record CreateStateRequest(
    @NotBlank @Size(min = 2, max = 2) String code,
    @NotBlank String name
) {
}
