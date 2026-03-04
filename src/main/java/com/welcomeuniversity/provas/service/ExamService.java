package com.welcomeuniversity.provas.service;

import java.time.Instant;
import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.server.ResponseStatusException;

import com.welcomeuniversity.provas.dto.exam.ExamResponse;
import com.welcomeuniversity.provas.dto.exam.ExamReviewRequest;
import com.welcomeuniversity.provas.dto.exam.ExamUploadRequest;
import com.welcomeuniversity.provas.model.AppUser;
import com.welcomeuniversity.provas.model.Exam;
import com.welcomeuniversity.provas.model.ExamStatus;
import com.welcomeuniversity.provas.model.Role;
import com.welcomeuniversity.provas.model.Subject;
import com.welcomeuniversity.provas.repository.ExamRepository;
import com.welcomeuniversity.provas.repository.SubjectRepository;

@Service
public class ExamService {

    private static final int MAX_PENDING_EXAMS_FOR_USER = 5;

    private final ExamRepository examRepository;
    private final SubjectRepository subjectRepository;
    private final CurrentUserService currentUserService;
    private final S3Service s3Service;

    public ExamService(
        ExamRepository examRepository,
        SubjectRepository subjectRepository,
        CurrentUserService currentUserService,
        S3Service s3Service
    ) {
        this.examRepository = examRepository;
        this.subjectRepository = subjectRepository;
        this.currentUserService = currentUserService;
        this.s3Service = s3Service;
    }

    public List<ExamResponse> listApproved(Long subjectId, String period) {
        if (period != null && !period.isBlank() && subjectId == null) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "O filtro period exige subjectId.");
        }

        if (subjectId == null) {
            return examRepository.findByStatusOrderByCreatedAtDesc(ExamStatus.APPROVED)
                .stream()
                .map(ExamResponse::from)
                .toList();
        }

        if (period == null || period.isBlank()) {
            return examRepository.findByStatusAndSubjectIdOrderByCreatedAtDesc(ExamStatus.APPROVED, subjectId)
                .stream()
                .map(ExamResponse::from)
                .toList();
        }

        int[] parsedPeriod = parsePeriod(period);
        return examRepository.findByStatusAndSubjectIdAndExamYearAndSemesterOrderByCreatedAtDesc(
                ExamStatus.APPROVED,
                subjectId,
                parsedPeriod[0],
                parsedPeriod[1]
            )
            .stream()
            .map(ExamResponse::from)
            .toList();
    }

    public List<ExamResponse> listPending() {
        return examRepository.findByStatusOrderByIdAsc(ExamStatus.PENDING)
            .stream()
            .map(ExamResponse::from)
            .toList();
    }

    @Transactional
    public ExamResponse upload(ExamUploadRequest request) {
        MultipartFile file = request.getFile();
        validatePdf(file);

        AppUser currentUser = currentUserService.requireCurrentUser();
        if (currentUser.getRole() == Role.USER) {
            long pendingCount = examRepository.countByUploadedByIdAndStatus(currentUser.getId(), ExamStatus.PENDING);
            if (pendingCount >= MAX_PENDING_EXAMS_FOR_USER) {
                throw new ResponseStatusException(
                    HttpStatus.CONFLICT,
                    "Usuarios com ROLE_USER podem ter no maximo 5 provas pendentes."
                );
            }
        }

        Subject subject = subjectRepository.findById(request.getSubjectId())
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Materia nao encontrada."));

        S3Service.StoredObject storedObject = s3Service.uploadExam(file, subject.getId());

        Exam exam = new Exam();
        exam.setName(request.getName().trim());
        exam.setExamYear(request.getExamYear());
        exam.setSemester(request.getSemester());
        exam.setType(request.getType());
        exam.setPdfUrl(storedObject.url());
        exam.setStorageKey(storedObject.key());
        exam.setStatus(ExamStatus.PENDING);
        exam.setSubject(subject);
        exam.setUploadedBy(currentUser);

        return ExamResponse.from(examRepository.save(exam));
    }

    @Transactional
    public ExamResponse review(Long examId, ExamReviewRequest request) {
        if (request.status() == ExamStatus.PENDING) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Nao e permitido voltar prova para PENDING.");
        }

        Exam exam = examRepository.findById(examId)
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Prova nao encontrada."));

        if (exam.getStatus() != ExamStatus.PENDING) {
            throw new ResponseStatusException(HttpStatus.CONFLICT, "Apenas provas pendentes podem ser validadas.");
        }

        AppUser reviewer = currentUserService.requireCurrentUser();
        exam.setStatus(request.status());
        exam.setReviewedBy(reviewer);
        exam.setReviewedAt(Instant.now());
        exam.setReviewNote(normalizeOptionalText(request.reviewNote()));

        return ExamResponse.from(exam);
    }

    private int[] parsePeriod(String period) {
        String[] parts = period.split("\\.");
        if (parts.length != 2) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Period deve seguir o formato AAAA.S.");
        }

        try {
            int examYear = Integer.parseInt(parts[0]);
            int semester = Integer.parseInt(parts[1]);
            if (semester < 1 || semester > 2) {
                throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Semestre invalido.");
            }
            return new int[] {examYear, semester};
        } catch (NumberFormatException ex) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Period deve conter apenas numeros.");
        }
    }

    private void validatePdf(MultipartFile file) {
        if (file == null || file.isEmpty()) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Arquivo PDF obrigatorio.");
        }

        String filename = file.getOriginalFilename() == null ? "" : file.getOriginalFilename().toLowerCase();
        String contentType = file.getContentType() == null ? "" : file.getContentType().toLowerCase();
        boolean validPdf = filename.endsWith(".pdf") || "application/pdf".equals(contentType);

        if (!validPdf) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Somente arquivos PDF sao aceitos.");
        }
    }

    private String normalizeOptionalText(String input) {
        if (input == null || input.isBlank()) {
            return null;
        }
        return input.trim();
    }
}
