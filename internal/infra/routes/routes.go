package routes

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
	"github.com/mikaelcaua/welcome-university-api/internal/infra/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/middleware"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/course"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/state"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/storage"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/university"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/repositories_impl/user"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/auth"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/course"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/state"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/university"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/user"
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
	storageRepository, err := storagerepository.NewAwsS3StorageRepositoryImpl(ctx, appConfig)
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
	router.GET("/actuator/health/readiness", health)
	router.GET("/actuator/health/liveness", health)
	router.GET("/actuator/info", func(ctx *gin.Context) { httpx.RespondJSON(ctx, http.StatusOK, gin.H{"app": "welcome-university-api"}) })
	authroutes.Register(router, authController)
	stateroutes.RegisterPublic(router, stateController)
	universityroutes.RegisterPublic(router, universityController)
	courseroutes.RegisterPublic(router, courseController)
	subjectroutes.RegisterPublic(router, subjectController)
	examroutes.RegisterPublic(router, examController)

	protected := router.Group("", middleware.Authentication(jwtService, userRepository))
	stateroutes.RegisterProtected(protected, stateController)
	universityroutes.RegisterProtected(protected, universityController)
	courseroutes.RegisterProtected(protected, courseController)
	subjectroutes.RegisterProtected(protected, subjectController)
	examroutes.RegisterProtected(protected, examController)
	userroutes.RegisterProtected(protected, userController)
	return router, func() {}, nil
}

func health(ctx *gin.Context) { httpx.RespondJSON(ctx, http.StatusOK, gin.H{"status": "UP"}) }
