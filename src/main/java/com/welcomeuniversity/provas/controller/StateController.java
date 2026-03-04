package com.welcomeuniversity.provas.controller;

import java.util.List;
import java.util.Locale;

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
import com.welcomeuniversity.provas.dto.state.CreateStateRequest;
import com.welcomeuniversity.provas.model.State;
import com.welcomeuniversity.provas.repository.StateRepository;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import jakarta.validation.Valid;

@RestController
@RequestMapping("/states")
public class StateController {

    private final StateRepository repo;

    public StateController(StateRepository repo){
        this.repo = repo;
    }

    @GetMapping
    @Operation(summary = "Listar estados", security = {})
    public List<State> list() {
        return repo.findAll();
    }

    @GetMapping("/{code}")
    @Operation(summary = "Buscar estado por sigla", security = {})
    public State get(@PathVariable String code) {
        return repo.findByCode(code)
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Estado não encontrado: " + code));
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(summary = "Criar estado", security = @SecurityRequirement(name = OpenApiConfig.BEARER_SCHEME))
    public State create(@Valid @RequestBody CreateStateRequest request) {
        String normalizedCode = request.code().trim().toUpperCase(Locale.ROOT);
        if (repo.findByCode(normalizedCode).isPresent()) {
            throw new ResponseStatusException(HttpStatus.CONFLICT, "Estado ja cadastrado: " + normalizedCode);
        }

        return repo.save(new State(normalizedCode, request.name().trim()));
    }
}
