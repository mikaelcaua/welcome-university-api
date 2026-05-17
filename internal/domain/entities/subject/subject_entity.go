package subject

import (
	"strings"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/course"
)

type Subject struct {
	ID     int64                `json:"id"`
	Name   string               `json:"name"`
	Course *course.Course `json:"course,omitempty"`
}

func NewSubject(name string) (Subject, bool) {
	normalizedName := strings.TrimSpace(name)
	if normalizedName == "" {
		return Subject{}, false
	}
	return Subject{Name: normalizedName}, true
}
