package courseusecase

import (
	"context"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/course"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/course"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/contracts/university"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/erros"
)

type CreateCourseUseCase struct {
	universities universitycontract.UniversityRepositoryContract
	courses      coursecontract.CourseRepositoryContract
}

func NewCreateCourseUseCase(universities universitycontract.UniversityRepositoryContract, courses coursecontract.CourseRepositoryContract) *CreateCourseUseCase {
	return &CreateCourseUseCase{universities: universities, courses: courses}
}
func (useCase *CreateCourseUseCase) Execute(ctx context.Context, universityID int64, name string) (course.Course, error) {
	_, found, err := useCase.universities.FindByID(ctx, universityID)
	if err != nil {
		return course.Course{}, err
	}
	if !found {
		return course.Course{}, domainerrors.NotFound("Universidade nao encontrada.")
	}
	course, valid := course.NewCourse(name)
	if !valid {
		return course.Course{}, domainerrors.Validation("Nome do curso e obrigatorio.")
	}
	exists, err := useCase.courses.NameExistsInUniversity(ctx, course.Name, universityID)
	if err != nil {
		return course.Course{}, err
	}
	if exists {
		return course.Course{}, domainerrors.Conflict("Curso ja cadastrado nesta universidade.")
	}
	return useCase.courses.Create(ctx, course, universityID)
}
