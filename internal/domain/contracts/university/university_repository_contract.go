package universitycontract

import (
	"context"

	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/university"
)

type UniversityRepositoryContract interface {
	ListByState(ctx context.Context, stateID int64) ([]university.University, error)
	FindByID(ctx context.Context, universityID int64) (university.University, bool, error)
	NameExistsInState(ctx context.Context, name string, stateID int64) (bool, error)
	Create(ctx context.Context, university university.University, stateID int64) (university.University, error)
}
