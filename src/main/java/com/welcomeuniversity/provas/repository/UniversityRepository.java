package com.welcomeuniversity.provas.repository;

import java.util.List;
import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;

import com.welcomeuniversity.provas.model.University;

public interface UniversityRepository extends JpaRepository<University, Long> {
    List<University> findByStateId(Long stateId);
    Optional<University> findByNameAndStateId(String name, Long stateId);
}
