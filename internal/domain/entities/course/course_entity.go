package course

import (
	"strings"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/university"
)

type Course struct {
	ID         int64                        `json:"id"`
	Name       string                       `json:"name"`
	University *university.University `json:"university,omitempty"`
}

func NewCourse(name string) (Course, bool) {
	normalizedName := strings.TrimSpace(name)
	if normalizedName == "" {
		return Course{}, false
	}
	return Course{Name: normalizedName}, true
}
