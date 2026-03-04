package com.welcomeuniversity.provas.controller;

import java.util.List;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.welcomeuniversity.provas.model.University;
import com.welcomeuniversity.provas.repository.UniversityRepository;

import io.swagger.v3.oas.annotations.Operation;

@RestController
@RequestMapping("/states/{stateId}/universities")
public class UniversityController {
    private final UniversityRepository repo;
    public UniversityController(UniversityRepository repo){ this.repo = repo; }

    @GetMapping
    @Operation(summary = "Listar universidades por estado", security = {})
    public List<University> list(@PathVariable Long stateId){ return repo.findByStateId(stateId); }
}
