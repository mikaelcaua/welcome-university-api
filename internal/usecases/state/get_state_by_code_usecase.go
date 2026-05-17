package stateusecase

import (
	"context"
	"strings"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/state"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/state"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type GetStateByCodeUseCase struct {
	states statecontract.StateRepositoryContract
}

func NewGetStateByCodeUseCase(states statecontract.StateRepositoryContract) *GetStateByCodeUseCase {
	return &GetStateByCodeUseCase{states: states}
}
func (useCase *GetStateByCodeUseCase) Execute(ctx context.Context, code string) (state.State, error) {
	foundState, found, err := useCase.states.FindByCode(ctx, strings.ToUpper(strings.TrimSpace(code)))
	if err != nil {
		return state.State{}, err
	}
	if !found {
		return state.State{}, domainerrors.NotFound("Estado não encontrado: " + code)
	}
	return foundState, nil
}
