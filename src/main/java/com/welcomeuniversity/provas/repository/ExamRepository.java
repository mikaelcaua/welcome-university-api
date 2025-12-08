package com.welcomeuniversity.provas.repository;

import java.util.List;

import org.springframework.data.jpa.repository.JpaRepository;

import com.welcomeuniversity.provas.model.Exam;

public interface ExamRepository extends JpaRepository<Exam, Long> {

    List<Exam> findBySubjectId(Long subjectId);

    List<Exam> findBySubjectIdAndExamYearAndSemester(Long subjectId, int examYear, int semester);
}
