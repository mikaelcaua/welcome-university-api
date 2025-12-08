package com.welcomeuniversity.provas.controller;

import java.util.List;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.welcomeuniversity.provas.model.Exam;
import com.welcomeuniversity.provas.repository.ExamRepository;

@RestController
@RequestMapping("/subjects/{subjectId}/exams")
public class ExamController {

    private final ExamRepository repo;

    public ExamController(ExamRepository repo){ 
        this.repo = repo; 
    }

    @GetMapping
    public List<Exam> listAll(
        @PathVariable Long subjectId, 
        @RequestParam(required = false) String period
    ) {
        if (period == null || period.isBlank()) {
            return repo.findBySubjectId(subjectId);
        }

        String[] parts = period.split("\\.");
        int examYear = Integer.parseInt(parts[0]);
        int semester = Integer.parseInt(parts[1]);

        return repo.findBySubjectIdAndExamYearAndSemester(subjectId, examYear, semester);
    }
}
