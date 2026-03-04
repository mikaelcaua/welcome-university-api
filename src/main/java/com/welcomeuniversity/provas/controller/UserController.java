package com.welcomeuniversity.provas.controller;

import java.util.List;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PatchMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.welcomeuniversity.provas.config.OpenApiConfig;
import com.welcomeuniversity.provas.dto.user.UpdateUserRoleRequest;
import com.welcomeuniversity.provas.dto.user.UserResponse;
import com.welcomeuniversity.provas.service.UserService;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import jakarta.validation.Valid;

@RestController
@RequestMapping("/users")
public class UserController {

    private final UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    @GetMapping("/me")
    @Operation(summary = "Obter usuario autenticado", security = @SecurityRequirement(name = OpenApiConfig.BEARER_SCHEME))
    public UserResponse me() {
        return userService.me();
    }

    @GetMapping
    @Operation(summary = "Listar usuarios", security = @SecurityRequirement(name = OpenApiConfig.BEARER_SCHEME))
    public List<UserResponse> listAll() {
        return userService.listAll();
    }

    @PatchMapping("/{id}/role")
    @Operation(summary = "Atualizar papel do usuario", security = @SecurityRequirement(name = OpenApiConfig.BEARER_SCHEME))
    public UserResponse updateRole(@PathVariable Long id, @Valid @RequestBody UpdateUserRoleRequest request) {
        return userService.updateRole(id, request.role());
    }
}
