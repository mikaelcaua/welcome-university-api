package com.welcomeuniversity.provas.controller;

import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;

import com.welcomeuniversity.provas.model.State;
import com.welcomeuniversity.provas.repository.StateRepository;

import io.swagger.v3.oas.annotations.Operation;

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
}
