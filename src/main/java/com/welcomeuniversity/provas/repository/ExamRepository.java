package com.welcomeuniversity.provas.repository;

import java.util.List;

import org.springframework.data.jpa.repository.JpaRepository;

import com.welcomeuniversity.provas.model.Exam;
import com.welcomeuniversity.provas.model.ExamStatus;

public interface ExamRepository extends JpaRepository<Exam, Long> {

    List<Exam> findByStatusOrderByCreatedAtDesc(ExamStatus status);

    List<Exam> findByStatusAndSubjectIdOrderByCreatedAtDesc(ExamStatus status, Long subjectId);

    List<Exam> findByStatusAndSubjectIdAndExamYearAndSemesterOrderByCreatedAtDesc(
        ExamStatus status,
        Long subjectId,
        int examYear,
        int semester
    );

    List<Exam> findByStatusOrderByIdAsc(ExamStatus status);
    List<Exam> findByStatusAndSubjectIdOrderByIdAsc(ExamStatus status, Long subjectId);

    long countByUploadedByIdAndStatus(Long uploadedById, ExamStatus status);

    boolean existsByFileHash(String fileHash);
}
