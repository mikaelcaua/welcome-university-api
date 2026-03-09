package com.welcomeuniversity.provas.dto.exam;

import java.time.Instant;

import com.welcomeuniversity.provas.dto.user.UserSummaryResponse;
import com.welcomeuniversity.provas.model.Exam;
import com.welcomeuniversity.provas.model.ExamStatus;
import com.welcomeuniversity.provas.model.ExamType;

public record ExamResponse(
    Long id,
    String name,
    int examYear,
    int semester,
    String periodLabel,
    ExamType type,
    String pdfUrl,
    ExamStatus status,
    Long subjectId,
    String subjectName,
    UserSummaryResponse uploadedBy,
    UserSummaryResponse reviewedBy,
    String reviewNote,
    Instant createdAt,
    Instant reviewedAt
) {

    public static ExamResponse from(Exam exam) {
        return new ExamResponse(
            exam.getId(),
            exam.getName(),
            exam.getExamYear(),
            exam.getSemester(),
            exam.getPeriodLabel(),
            exam.getType(),
            exam.getPdfUrl(),
            exam.getStatus(),
            exam.getSubject() != null ? exam.getSubject().getId() : null,
            exam.getSubject() != null ? exam.getSubject().getName() : null,
            UserSummaryResponse.from(exam.getUploadedBy()),
            UserSummaryResponse.from(exam.getReviewedBy()),
            exam.getReviewNote(),
            exam.getCreatedAt(),
            exam.getReviewedAt()
        );
    }
}
