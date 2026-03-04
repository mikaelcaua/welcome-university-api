package com.welcomeuniversity.provas.controller;

import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PatchMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

import com.welcomeuniversity.provas.dto.exam.ExamResponse;
import com.welcomeuniversity.provas.dto.exam.ExamReviewRequest;
import com.welcomeuniversity.provas.dto.exam.ExamUploadRequest;
import com.welcomeuniversity.provas.config.OpenApiConfig;
import com.welcomeuniversity.provas.service.ExamService;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import jakarta.validation.Valid;

@RestController
public class ExamController {

    private final ExamService examService;

    public ExamController(ExamService examService) {
        this.examService = examService;
    }

    @GetMapping("/subjects/{subjectId}/exams")
    @Operation(summary = "Listar provas aprovadas por disciplina", security = {})
    public List<ExamResponse> listBySubject(
        @PathVariable Long subjectId,
        @RequestParam(required = false) String period
    ) {
        return examService.listApproved(subjectId, period);
    }

    @GetMapping("/exams")
    @Operation(summary = "Listar provas aprovadas", security = {})
    public List<ExamResponse> listAll(
        @RequestParam(required = false) Long subjectId,
        @RequestParam(required = false) String period
    ) {
        return examService.listApproved(subjectId, period);
    }

    @GetMapping("/exams/pending")
    @Operation(summary = "Listar provas pendentes", security = @SecurityRequirement(name = OpenApiConfig.BEARER_SCHEME))
    public List<ExamResponse> listPending() {
        return examService.listPending();
    }

    @PostMapping(value = "/exams", consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(summary = "Enviar prova", security = @SecurityRequirement(name = OpenApiConfig.BEARER_SCHEME))
    public ExamResponse upload(@Valid @ModelAttribute ExamUploadRequest request) {
        return examService.upload(request);
    }

    @PatchMapping("/exams/{examId}/status")
    @Operation(summary = "Revisar status da prova", security = @SecurityRequirement(name = OpenApiConfig.BEARER_SCHEME))
    public ExamResponse review(@PathVariable Long examId, @Valid @RequestBody ExamReviewRequest request) {
        return examService.review(examId, request);
    }
}
