package com.welcomeuniversity.provas.repository;

import java.util.List;

import org.springframework.data.jpa.repository.JpaRepository;

import com.welcomeuniversity.provas.model.Subject;

public interface SubjectRepository extends JpaRepository<Subject, Long> {
    List<Subject> findByCourseId(Long courseId);
}
