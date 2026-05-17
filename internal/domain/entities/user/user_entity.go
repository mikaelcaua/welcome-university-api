package user

import (
	"strings"
	"time"
)

type Role string

const (
	RoleUser     Role = "USER"
	RoleApprover Role = "APPROVER"
	RoleAdmin    Role = "ADMIN"
	RoleDev      Role = "DEV"
)

func (role Role) IsValidAssignableRole() bool {
	return role == RoleUser || role == RoleApprover || role == RoleAdmin
}

type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (user User) HasValidIdentity() bool {
	return strings.TrimSpace(user.Name) != "" && strings.TrimSpace(user.Email) != ""
}
