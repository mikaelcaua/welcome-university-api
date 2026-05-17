package universityusecase

import (
	"context"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/university"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/state"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/university"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type CreateUniversityUseCase struct {
	states       statecontract.StateRepositoryContract
	universities universitycontract.UniversityRepositoryContract
}

func NewCreateUniversityUseCase(states statecontract.StateRepositoryContract, universities universitycontract.UniversityRepositoryContract) *CreateUniversityUseCase {
	return &CreateUniversityUseCase{states: states, universities: universities}
}
func (useCase *CreateUniversityUseCase) Execute(ctx context.Context, stateID int64, name string) (university.University, error) {
	_, found, err := useCase.states.FindByID(ctx, stateID)
	if err != nil {
		return university.University{}, err
	}
	if !found {
		return university.University{}, domainerrors.NotFound("Estado nao encontrado.")
	}
	university, valid := university.NewUniversity(name)
	if !valid {
		return university.University{}, domainerrors.Validation("Nome da universidade e obrigatorio.")
	}
	exists, err := useCase.universities.NameExistsInState(ctx, university.Name, stateID)
	if err != nil {
		return university.University{}, err
	}
	if exists {
		return university.University{}, domainerrors.Conflict("Universidade ja cadastrada neste estado.")
	}
	return useCase.universities.Create(ctx, university, stateID)
}
