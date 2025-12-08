package com.welcomeuniversity.provas.repository;

import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;

import com.welcomeuniversity.provas.model.State;

public interface StateRepository extends JpaRepository<State, Long> {
    Optional<State> findByCode(String code);
}