package com.welcomeuniversity.provas.controller;

import java.util.List;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.welcomeuniversity.provas.model.State;
import com.welcomeuniversity.provas.repository.StateRepository;

@RestController
@RequestMapping("/states")
public class StateController {

    private final StateRepository repo;
    public StateController(StateRepository repo){ this.repo = repo; }

    @GetMapping
    public List<State> list() { return repo.findAll(); }

    @GetMapping("/{id}")
    public State get(@PathVariable Long id){ return repo.findById(id).orElseThrow(); }
}
