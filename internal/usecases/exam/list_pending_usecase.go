package examusecase

import (
	"context"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type ListPendingUseCase struct {
	exams    examcontract.ExamRepositoryContract
	subjects subjectcontract.SubjectRepositoryContract
}

func NewListPendingUseCase(exams examcontract.ExamRepositoryContract, subjects subjectcontract.SubjectRepositoryContract) *ListPendingUseCase {
	return &ListPendingUseCase{exams: exams, subjects: subjects}
}
func (useCase *ListPendingUseCase) Execute(ctx context.Context, stateID int64, universityID int64, courseID int64, subjectID int64) ([]exam.Exam, error) {
	subject, found, err := useCase.subjects.FindTreeByID(ctx, subjectID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, domainerrors.NotFound("Materia nao encontrada.")
	}
	if subject.Course == nil || subject.Course.ID != courseID {
		return nil, domainerrors.Validation("subjectId nao pertence ao courseId informado.")
	}
	if subject.Course.University == nil || subject.Course.University.ID != universityID {
		return nil, domainerrors.Validation("courseId nao pertence ao universityId informado.")
	}
	if subject.Course.University.State == nil || subject.Course.University.State.ID != stateID {
		return nil, domainerrors.Validation("universityId nao pertence ao stateId informado.")
	}
	return useCase.exams.ListPendingBySubject(ctx, subjectID)
}
