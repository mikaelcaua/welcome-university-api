package stateusecase

import (
	"context"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/state"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/state"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type CreateStateUseCase struct{ states statecontract.StateRepositoryContract }

func NewCreateStateUseCase(states statecontract.StateRepositoryContract) *CreateStateUseCase {
	return &CreateStateUseCase{states: states}
}
func (useCase *CreateStateUseCase) Execute(ctx context.Context, code string, name string) (state.State, error) {
	state, valid := state.NewState(code, name)
	if !valid {
		return state.State{}, domainerrors.Validation("Sigla e nome sao obrigatorios.")
	}
	_, found, err := useCase.states.FindByCode(ctx, state.Code)
	if err != nil {
		return state.State{}, err
	}
	if found {
		return state.State{}, domainerrors.Conflict("Estado ja cadastrado: "+state.Code)
	}
	return useCase.states.Create(ctx, state)
}
