package bootstrap

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/config"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/auth"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/course"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/state"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/university"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/user"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/middleware"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/course"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/state"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/storage"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/university"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/user"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/security"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/auth"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/course"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/state"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/university"
	"github.com/mikaelcaua/welcome-university-api/internal/usecases/user"
)

func Create(ctx context.Context, appConfig config.Config, pool *pgxpool.Pool) (http.Handler, func(), error) {
	userRepository := userrepository.NewPostgresUserRepositoryImpl(pool)
	stateRepository := staterepository.NewPostgresStateRepositoryImpl(pool)
	universityRepository := universityrepository.NewPostgresUniversityRepositoryImpl(pool)
	courseRepository := courserepository.NewPostgresCourseRepositoryImpl(pool)
	subjectRepository := subjectrepository.NewPostgresSubjectRepositoryImpl(pool)
	examRepository := examrepository.NewPostgresExamRepositoryImpl(pool)
	storageRepository, err := storagerepository.NewMinioStorageRepositoryImpl(ctx, appConfig)
	if err != nil {
		return nil, nil, err
	}

	jwtService := security.NewJWTService(appConfig.JWTSecret, appConfig.AccessTokenExpirationSeconds, appConfig.RefreshTokenExpirationSeconds)
	authController := authcontroller.NewAuthController(authusecase.NewRegisterUseCase(userRepository, jwtService, appConfig.AccessTokenExpirationSeconds), authusecase.NewLoginUseCase(userRepository, jwtService, appConfig.AccessTokenExpirationSeconds), authusecase.NewRefreshUseCase(userRepository, jwtService, appConfig.AccessTokenExpirationSeconds))
	userController := usercontroller.NewUserController(userusecase.NewGetCurrentUserUseCase(), userusecase.NewListUsersUseCase(userRepository), userusecase.NewUpdateUserRoleUseCase(userRepository))
	stateController := statecontroller.NewStateController(stateusecase.NewListStatesUseCase(stateRepository), stateusecase.NewGetStateByCodeUseCase(stateRepository), stateusecase.NewCreateStateUseCase(stateRepository))
	universityController := universitycontroller.NewUniversityController(universityusecase.NewListUniversitiesUseCase(universityRepository), universityusecase.NewCreateUniversityUseCase(stateRepository, universityRepository))
	courseController := coursecontroller.NewCourseController(courseusecase.NewListCoursesUseCase(courseRepository), courseusecase.NewCreateCourseUseCase(universityRepository, courseRepository))
	subjectController := subjectcontroller.NewSubjectController(subjectusecase.NewListSubjectsUseCase(subjectRepository), subjectusecase.NewCreateSubjectUseCase(courseRepository, subjectRepository))
	examController := examcontroller.NewExamController(examusecase.NewListApprovedUseCase(examRepository), examusecase.NewListPendingUseCase(examRepository, subjectRepository), examusecase.NewListMyPendingUseCase(examRepository), examusecase.NewUploadExamUseCase(examRepository, subjectRepository, storageRepository, examusecase.NewUploadOptimizer()), examusecase.NewReviewExamUseCase(examRepository, storageRepository), appConfig.MaxUploadSizeMB)

	router := gin.New()
	router.Use(gin.Recovery())
	routes.Register(router, routes.Dependencies{
		AuthController:       authController,
		UserController:       userController,
		StateController:      stateController,
		UniversityController: universityController,
		CourseController:     courseController,
		SubjectController:    subjectController,
		ExamController:       examController,
		Authentication:       middleware.Authentication(jwtService, userRepository),
	})
	return router, func() {}, nil
}
