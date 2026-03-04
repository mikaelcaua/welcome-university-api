package com.welcomeuniversity.provas.repository;

import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;

import com.welcomeuniversity.provas.model.AppUser;

public interface UserRepository extends JpaRepository<AppUser, Long> {

    boolean existsByEmailIgnoreCase(String email);

    Optional<AppUser> findByEmailIgnoreCase(String email);
}
