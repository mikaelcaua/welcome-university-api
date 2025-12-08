package com.welcomeuniversity.provas.controller;

import java.util.List;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.welcomeuniversity.provas.model.Subject;
import com.welcomeuniversity.provas.repository.SubjectRepository;

@RestController
@RequestMapping("/courses/{courseId}/subjects")
public class SubjectController {
    private final SubjectRepository repo;
    public SubjectController(SubjectRepository repo){ this.repo = repo; }

    @GetMapping
    public List<Subject> list(@PathVariable Long courseId){ return repo.findByCourseId(courseId); }
}
