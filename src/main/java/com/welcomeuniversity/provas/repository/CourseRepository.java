package com.welcomeuniversity.provas.repository;

import java.util.List;

import org.springframework.data.jpa.repository.JpaRepository;

import com.welcomeuniversity.provas.model.Course;

public interface CourseRepository extends JpaRepository<Course, Long> {
    List<Course> findByUniversityId(Long universityId);
}
