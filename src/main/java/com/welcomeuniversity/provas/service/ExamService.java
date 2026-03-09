package com.welcomeuniversity.provas.service;

import java.time.Instant;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.List;
import java.util.Locale;
import java.util.Set;

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
    private static final String PERIOD_LABEL_UNIDENTIFIED = "PERIODO_NAO_IDENTIFICADO";
    private static final Set<String> SUPPORTED_FILE_EXTENSIONS = Set.of(
        "pdf",
        "png",
        "jpg",
        "jpeg",
        "webp",
        "gif"
    );

    private final ExamRepository examRepository;
    private final SubjectRepository subjectRepository;
    private final CurrentUserService currentUserService;
    private final S3Service s3Service;
    private final UploadCompressionService uploadCompressionService;

    public ExamService(
        ExamRepository examRepository,
        SubjectRepository subjectRepository,
        CurrentUserService currentUserService,
        S3Service s3Service,
        UploadCompressionService uploadCompressionService
    ) {
        this.examRepository = examRepository;
        this.subjectRepository = subjectRepository;
        this.currentUserService = currentUserService;
        this.s3Service = s3Service;
        this.uploadCompressionService = uploadCompressionService;
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

    public List<ExamResponse> listPending(Long stateId, Long universityId, Long courseId, Long subjectId) {
        Subject subject = subjectRepository.findById(subjectId)
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Materia nao encontrada."));

        if (subject.getCourse() == null || subject.getCourse().getId() == null || !subject.getCourse().getId().equals(courseId)) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "subjectId nao pertence ao courseId informado.");
        }
        if (subject.getCourse().getUniversity() == null
            || subject.getCourse().getUniversity().getId() == null
            || !subject.getCourse().getUniversity().getId().equals(universityId)) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "courseId nao pertence ao universityId informado.");
        }
        if (subject.getCourse().getUniversity().getState() == null
            || subject.getCourse().getUniversity().getState().getId() == null
            || !subject.getCourse().getUniversity().getState().getId().equals(stateId)) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "universityId nao pertence ao stateId informado.");
        }

        return examRepository.findByStatusAndSubjectIdOrderByIdAsc(ExamStatus.PENDING, subjectId)
            .stream()
            .map(ExamResponse::from)
            .toList();
    }

    @Transactional
    public ExamResponse upload(ExamUploadRequest request) {
        MultipartFile file = request.getFile();
        validateSupportedFile(file);

        AppUser currentUser = currentUserService.requireCurrentUser();
        if (currentUser.getRole() != Role.ADMIN && currentUser.getRole() != Role.DEV) {
            long pendingCount = examRepository.countByUploadedByIdAndStatus(currentUser.getId(), ExamStatus.PENDING);
            if (pendingCount >= MAX_PENDING_EXAMS_FOR_USER) {
                throw new ResponseStatusException(
                    HttpStatus.CONFLICT,
                    "Somente ADMIN/DEV podem ultrapassar o limite de 5 provas pendentes."
                );
            }
        }

        Subject subject = subjectRepository.findById(request.getSubjectId())
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Materia nao encontrada."));

        S3Service.UploadPayload payload = uploadCompressionService.compressForUpload(file);
        String fileHash = sha256Hex(payload.bytes());
        if (examRepository.existsByFileHash(fileHash)) {
            throw new ResponseStatusException(HttpStatus.CONFLICT, "Arquivo ja existe na base.");
        }

        S3Service.StoredObject storedObject = s3Service.uploadExam(payload, subject.getId());
        boolean periodUnidentified = Boolean.TRUE.equals(request.getPeriodUnidentified());

        Exam exam = new Exam();
        exam.setName(buildExamName(subject, request, periodUnidentified));
        exam.setExamYear(request.getExamYear());
        exam.setSemester(request.getSemester());
        exam.setPeriodLabel(periodUnidentified ? PERIOD_LABEL_UNIDENTIFIED : null);
        exam.setType(request.getType());
        exam.setPdfUrl(storedObject.url());
        exam.setStorageKey(storedObject.key());
        exam.setFileHash(fileHash);
        if (currentUser.getRole() == Role.ADMIN || currentUser.getRole() == Role.DEV) {
            exam.setStatus(ExamStatus.APPROVED);
            exam.setReviewedBy(currentUser);
            exam.setReviewedAt(Instant.now());
            exam.setReviewNote("Aprovacao automatica: upload por ADMIN/DEV.");
        } else {
            exam.setStatus(ExamStatus.PENDING);
        }
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

    private void validateSupportedFile(MultipartFile file) {
        if (file == null || file.isEmpty()) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Arquivo obrigatorio.");
        }

        String filename = file.getOriginalFilename() == null ? "" : file.getOriginalFilename().toLowerCase(Locale.ROOT);
        String contentType = file.getContentType() == null ? "" : file.getContentType().toLowerCase();
        String extension = extractFileExtension(filename);
        boolean validFile = "application/pdf".equals(contentType)
            || contentType.startsWith("image/")
            || SUPPORTED_FILE_EXTENSIONS.contains(extension);

        if (!validFile) {
            throw new ResponseStatusException(
                HttpStatus.BAD_REQUEST,
                "Somente arquivos PDF ou imagens sao aceitos."
            );
        }
    }

    private String buildExamName(Subject subject, ExamUploadRequest request, boolean periodUnidentified) {
        if (periodUnidentified) {
            return "%s - %s - %s".formatted(
                subject.getName().trim(),
                PERIOD_LABEL_UNIDENTIFIED,
                formatExamType(request.getType().name())
            );
        }

        return "%s - %d.%d - %s".formatted(
            subject.getName().trim(),
            request.getExamYear(),
            request.getSemester(),
            formatExamType(request.getType().name())
        );
    }

    private String formatExamType(String type) {
        return type.replace('_', ' ').toUpperCase(Locale.ROOT);
    }

    private String extractFileExtension(String filename) {
        int extensionStart = filename.lastIndexOf('.');
        if (extensionStart < 0 || extensionStart == filename.length() - 1) {
            return "";
        }
        return filename.substring(extensionStart + 1);
    }

    private String normalizeOptionalText(String input) {
        if (input == null || input.isBlank()) {
            return null;
        }
        return input.trim();
    }

    private String sha256Hex(byte[] content) {
        try {
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hash = digest.digest(content);
            StringBuilder builder = new StringBuilder(hash.length * 2);
            for (byte b : hash) {
                builder.append(String.format("%02x", b));
            }
            return builder.toString();
        } catch (NoSuchAlgorithmException ex) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Falha ao calcular hash do arquivo.");
        }
    }
}
