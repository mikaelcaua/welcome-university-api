package university

import (
	"strings"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/state"
)

type University struct {
	ID    int64              `json:"id"`
	Name  string             `json:"name"`
	State *state.State `json:"state,omitempty"`
}

func NewUniversity(name string) (University, bool) {
	normalizedName := strings.TrimSpace(name)
	if normalizedName == "" {
		return University{}, false
	}
	return University{Name: normalizedName}, true
}
