package statecontract

import (
	"context"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/state"
)

type StateRepositoryContract interface {
	List(ctx context.Context) ([]state.State, error)
	FindByID(ctx context.Context, stateID int64) (state.State, bool, error)
	FindByCode(ctx context.Context, code string) (state.State, bool, error)
	Create(ctx context.Context, state state.State) (state.State, error)
}
