package stateusecase

import (
	"context"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/state"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/state"
)

type ListStatesUseCase struct{ states statecontract.StateRepositoryContract }

func NewListStatesUseCase(states statecontract.StateRepositoryContract) *ListStatesUseCase {
	return &ListStatesUseCase{states: states}
}
func (useCase *ListStatesUseCase) Execute(ctx context.Context) ([]state.State, error) {
	return useCase.states.List(ctx)
}
