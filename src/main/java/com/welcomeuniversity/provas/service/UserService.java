package com.welcomeuniversity.provas.service;

import java.util.List;

import org.springframework.data.domain.Sort;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import com.welcomeuniversity.provas.dto.user.UserResponse;
import com.welcomeuniversity.provas.model.AppUser;
import com.welcomeuniversity.provas.model.Role;
import com.welcomeuniversity.provas.repository.UserRepository;

@Service
public class UserService {

    private final UserRepository userRepository;
    private final CurrentUserService currentUserService;

    public UserService(UserRepository userRepository, CurrentUserService currentUserService) {
        this.userRepository = userRepository;
        this.currentUserService = currentUserService;
    }

    public UserResponse me() {
        return UserResponse.from(currentUserService.requireCurrentUser());
    }

    public List<UserResponse> listAll() {
        return userRepository.findAll(Sort.by(Sort.Direction.ASC, "id"))
            .stream()
            .map(UserResponse::from)
            .toList();
    }

    @Transactional
    public UserResponse updateRole(Long userId, Role role) {
        if (role == Role.DEV) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "ROLE_DEV nao pode ser atribuido via API.");
        }

        AppUser user = userRepository.findById(userId)
            .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Usuario nao encontrado."));

        user.setRole(role);
        return UserResponse.from(user);
    }
}
