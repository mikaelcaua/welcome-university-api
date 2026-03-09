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
import com.welcomeuniversity.provas.dto.subject.CreateSubjectRequest;
import com.welcomeuniversity.provas.model.Course;
import com.welcomeuniversity.provas.model.Subject;
import com.welcomeuniversity.provas.repository.CourseRepository;
import com.welcomeuniversity.provas.repository.SubjectRepository;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import jakarta.validation.Valid;

@RestController
@RequestMapping("/courses/{courseId}/subjects")
public class SubjectController {
    private final SubjectRepository repo;
    private final CourseRepository courseRepository;
    public SubjectController(SubjectRepository repo, CourseRepository courseRepository){
        this.repo = repo;
        this.courseRepository = courseRepository;
    }

    @GetMapping
    @Operation(summary = "Listar disciplinas por curso", security = {})
    public List<Subject> list(@PathVariable Long courseId){ return repo.findByCourseId(courseId); }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(
        summary = "Criar disciplina por curso",
        description = "Disponivel apenas para perfis ADMIN ou DEV. Retorna 404 se o curso nao existir e 409 em duplicidade.",
        security = @SecurityRequirement(name = OpenApiConfig.BEARER_SCHEME)
    )
    public Subject create(@PathVariable Long courseId, @Valid @RequestBody CreateSubjectRequest request) {
        Course course = courseRepository.findById(courseId)
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Curso nao encontrado."));

        String normalizedName = request.name().trim();
        if (repo.findByNameAndCourseId(normalizedName, courseId).isPresent()) {
            throw new ResponseStatusException(HttpStatus.CONFLICT, "Disciplina ja cadastrada neste curso.");
        }

        return repo.save(new Subject(normalizedName, course));
    }
}
