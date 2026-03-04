package com.welcomeuniversity.provas.controller;

import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;

import com.welcomeuniversity.provas.config.OpenApiConfig;
import com.welcomeuniversity.provas.dto.course.CreateCourseRequest;
import com.welcomeuniversity.provas.model.Course;
import com.welcomeuniversity.provas.model.University;
import com.welcomeuniversity.provas.repository.CourseRepository;
import com.welcomeuniversity.provas.repository.UniversityRepository;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import jakarta.validation.Valid;

@RestController
@RequestMapping("/universities/{universityId}/courses")
public class CourseController {
    private final CourseRepository repo;
    private final UniversityRepository universityRepository;

    public CourseController(CourseRepository repo, UniversityRepository universityRepository){
        this.repo = repo;
        this.universityRepository = universityRepository;
    }

    @GetMapping
    @Operation(summary = "Listar cursos por universidade", security = {})
    public List<Course> list(@PathVariable Long universityId){ return repo.findByUniversityId(universityId); }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(summary = "Criar curso por universidade", security = @SecurityRequirement(name = OpenApiConfig.BEARER_SCHEME))
    public Course create(@PathVariable Long universityId, @Valid @RequestBody CreateCourseRequest request) {
        University university = universityRepository.findById(universityId)
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Universidade nao encontrada."));

        String normalizedName = request.name().trim();
        if (repo.findByNameAndUniversityId(normalizedName, universityId).isPresent()) {
            throw new ResponseStatusException(HttpStatus.CONFLICT, "Curso ja cadastrado nesta universidade.");
        }

        return repo.save(new Course(normalizedName, university));
    }
}
