package subjectusecase

import (
	"context"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/course"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type CreateSubjectUseCase struct {
	courses  coursecontract.CourseRepositoryContract
	subjects subjectcontract.SubjectRepositoryContract
}

func NewCreateSubjectUseCase(courses coursecontract.CourseRepositoryContract, subjects subjectcontract.SubjectRepositoryContract) *CreateSubjectUseCase {
	return &CreateSubjectUseCase{courses: courses, subjects: subjects}
}
func (useCase *CreateSubjectUseCase) Execute(ctx context.Context, courseID int64, name string) (subject.Subject, error) {
	_, found, err := useCase.courses.FindByID(ctx, courseID)
	if err != nil {
		return subject.Subject{}, err
	}
	if !found {
		return subject.Subject{}, domainerrors.NotFound("Curso nao encontrado.")
	}
	createdSubject, valid := subject.NewSubject(name)
	if !valid {
		return subject.Subject{}, domainerrors.Validation("Nome da disciplina e obrigatorio.")
	}
	exists, err := useCase.subjects.NameExistsInCourse(ctx, createdSubject.Name, courseID)
	if err != nil {
		return subject.Subject{}, err
	}
	if exists {
		return subject.Subject{}, domainerrors.Conflict("Disciplina ja cadastrada neste curso.")
	}
	return useCase.subjects.Create(ctx, createdSubject, courseID)
}
