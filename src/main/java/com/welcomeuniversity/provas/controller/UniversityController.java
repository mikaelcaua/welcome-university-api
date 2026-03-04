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
import com.welcomeuniversity.provas.dto.university.CreateUniversityRequest;
import com.welcomeuniversity.provas.model.State;
import com.welcomeuniversity.provas.model.University;
import com.welcomeuniversity.provas.repository.StateRepository;
import com.welcomeuniversity.provas.repository.UniversityRepository;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import jakarta.validation.Valid;

@RestController
@RequestMapping("/states/{stateId}/universities")
public class UniversityController {
    private final UniversityRepository repo;
    private final StateRepository stateRepository;

    public UniversityController(UniversityRepository repo, StateRepository stateRepository){
        this.repo = repo;
        this.stateRepository = stateRepository;
    }

    @GetMapping
    @Operation(summary = "Listar universidades por estado", security = {})
    public List<University> list(@PathVariable Long stateId){ return repo.findByStateId(stateId); }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(summary = "Criar universidade por estado", security = @SecurityRequirement(name = OpenApiConfig.BEARER_SCHEME))
    public University create(@PathVariable Long stateId, @Valid @RequestBody CreateUniversityRequest request) {
        State state = stateRepository.findById(stateId)
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Estado nao encontrado."));

        String normalizedName = request.name().trim();
        if (repo.findByNameAndStateId(normalizedName, stateId).isPresent()) {
            throw new ResponseStatusException(HttpStatus.CONFLICT, "Universidade ja cadastrada neste estado.");
        }

        return repo.save(new University(normalizedName, state));
    }
}
