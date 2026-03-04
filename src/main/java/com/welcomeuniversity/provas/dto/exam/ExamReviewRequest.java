package com.welcomeuniversity.provas.dto.exam;

import com.welcomeuniversity.provas.model.ExamStatus;

import jakarta.validation.constraints.NotNull;

public record ExamReviewRequest(
    @NotNull ExamStatus status,
    String reviewNote
) {
}
