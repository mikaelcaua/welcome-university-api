package models

import "time"

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

type ExamType string

const (
	ExamTypeFirst      ExamType = "PROVA1"
	ExamTypeSecond     ExamType = "PROVA2"
	ExamTypeThird      ExamType = "PROVA3"
	ExamTypeRecovery   ExamType = "RECUPERACAO"
	ExamTypeFinal      ExamType = "FINAL"
)

func (examType ExamType) IsValid() bool {
	switch examType {
	case ExamTypeFirst, ExamTypeSecond, ExamTypeThird, ExamTypeRecovery, ExamTypeFinal:
		return true
	default:
		return false
	}
}

type ExamStatus string

const (
	ExamStatusPending  ExamStatus = "PENDING"
	ExamStatusApproved ExamStatus = "APPROVED"
	ExamStatusRejected ExamStatus = "REJECTED"
)

func (status ExamStatus) IsReviewResult() bool {
	return status == ExamStatusApproved || status == ExamStatusRejected
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

type State struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type University struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	State *State `json:"state,omitempty"`
}

type Course struct {
	ID         int64       `json:"id"`
	Name       string      `json:"name"`
	University *University `json:"university,omitempty"`
}

type Subject struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Course *Course `json:"course,omitempty"`
}

type Exam struct {
	ID          int64
	Name        string
	ExamYear    int
	Semester    int
	PeriodLabel *string
	Type        ExamType
	PDFURL      string
	StorageKey  *string
	FileHash    *string
	Status      ExamStatus
	Subject     *Subject
	UploadedBy  *User
	ReviewedBy  *User
	ReviewNote  *string
	CreatedAt   time.Time
	ReviewedAt  *time.Time
	UpdatedAt   time.Time
}
