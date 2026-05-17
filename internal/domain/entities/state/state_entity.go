package state

import "strings"

type State struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func NewState(code string, name string) (State, bool) {
	normalizedCode := strings.ToUpper(strings.TrimSpace(code))
	normalizedName := strings.TrimSpace(name)
	if len(normalizedCode) != 2 || normalizedName == "" {
		return State{}, false
	}
	return State{Code: normalizedCode, Name: normalizedName}, true
}
