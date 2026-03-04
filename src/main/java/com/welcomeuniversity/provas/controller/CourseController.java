package com.welcomeuniversity.provas.controller;

import java.util.List;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.welcomeuniversity.provas.model.Course;
import com.welcomeuniversity.provas.repository.CourseRepository;

import io.swagger.v3.oas.annotations.Operation;

@RestController
@RequestMapping("/universities/{universityId}/courses")
public class CourseController {
    private final CourseRepository repo;
    public CourseController(CourseRepository repo){ this.repo = repo; }

    @GetMapping
    @Operation(summary = "Listar cursos por universidade", security = {})
    public List<Course> list(@PathVariable Long universityId){ return repo.findByUniversityId(universityId); }
}
