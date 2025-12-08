package com.welcomeuniversity.provas.repository;

import org.springframework.data.jpa.repository.JpaRepository;

import com.welcomeuniversity.provas.model.State;

public interface StateRepository extends JpaRepository<State, Long> {
}
