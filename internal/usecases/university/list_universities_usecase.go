package universityusecase

import (
	"context"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/university"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/university"
)

type ListUniversitiesUseCase struct{ universities universitycontract.UniversityRepositoryContract }

func NewListUniversitiesUseCase(universities universitycontract.UniversityRepositoryContract) *ListUniversitiesUseCase {
	return &ListUniversitiesUseCase{universities: universities}
}
func (useCase *ListUniversitiesUseCase) Execute(ctx context.Context, stateID int64) ([]university.University, error) {
	return useCase.universities.ListByState(ctx, stateID)
}
