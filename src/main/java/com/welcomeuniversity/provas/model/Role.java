package com.welcomeuniversity.provas.model;

public enum Role {
    USER,
    APPROVER,
    ADMIN,
    DEV;

    public String authority() {
        return "ROLE_" + name();
    }
}
