package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mikaelcaua/welcome-university-api/internal/config"
	"github.com/mikaelcaua/welcome-university-api/internal/controllers"
	"github.com/mikaelcaua/welcome-university-api/internal/httpx"
	authmiddleware "github.com/mikaelcaua/welcome-university-api/internal/middleware"
	"github.com/mikaelcaua/welcome-university-api/internal/models"
	"github.com/mikaelcaua/welcome-university-api/internal/repositories"
	"github.com/mikaelcaua/welcome-university-api/internal/security"
	"github.com/mikaelcaua/welcome-university-api/internal/services"
)

func Create(ctx context.Context, appConfig config.Config, pool *pgxpool.Pool) (http.Handler, func(), error) {
	userRepository := repositories.NewPostgresUserRepository(pool)
	catalogRepository := repositories.NewPostgresCatalogRepository(pool)
	examRepository := repositories.NewPostgresExamRepository(pool)

	jwtService := security.NewJWTService(appConfig.JWTSecret, appConfig.AccessTokenExpirationSeconds, appConfig.RefreshTokenExpirationSeconds)
	storageService, err := services.NewS3StorageService(ctx, appConfig)
	if err != nil {
		return nil, nil, err
	}
	uploadOptimizer := services.NewUploadOptimizer()

	authService := services.NewAuthService(userRepository, jwtService, appConfig.AccessTokenExpirationSeconds)
	userService := services.NewUserService(userRepository)
	catalogService := services.NewCatalogService(catalogRepository)
	examService := services.NewExamService(examRepository, catalogRepository, storageService, uploadOptimizer)

	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	catalogController := controllers.NewCatalogController(catalogService)
	examController := controllers.NewExamController(examService, appConfig.MaxUploadSizeMB)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	router.Get("/actuator/health/readiness", health)
	router.Get("/actuator/health/liveness", health)
	router.Get("/actuator/info", func(w http.ResponseWriter, r *http.Request) {
		httpx.RespondJSON(w, http.StatusOK, map[string]string{"app": "welcome-university-api"})
	})

	authRateLimiter := authmiddleware.NewRateLimiter(20, time.Minute)
	router.Group(func(authRoutes chi.Router) {
		authRoutes.Use(authRateLimiter.Middleware)
		authRoutes.Post("/auth/register", authController.Register)
		authRoutes.Post("/auth/login", authController.Login)
		authRoutes.Post("/auth/refresh", authController.Refresh)
	})

	router.Get("/states", catalogController.ListStates)
	router.Get("/states/{code}", catalogController.GetStateByCode)
	router.Get("/states/{stateId}/universities", catalogController.ListUniversities)
	router.Get("/universities/{universityId}/courses", catalogController.ListCourses)
	router.Get("/courses/{courseId}/subjects", catalogController.ListSubjects)
	router.Get("/subjects/{subjectId}/exams", examController.ListBySubject)
	router.Get("/exams", examController.ListAll)

	authenticated := authmiddleware.Authentication(jwtService, userRepository)
	router.Group(func(protected chi.Router) {
		protected.Use(authenticated)

		protected.With(authmiddleware.RequireRoles(models.RoleAdmin, models.RoleDev)).Post("/states", catalogController.CreateState)
		protected.With(authmiddleware.RequireRoles(models.RoleAdmin, models.RoleDev)).Post("/states/{stateId}/universities", catalogController.CreateUniversity)
		protected.With(authmiddleware.RequireRoles(models.RoleAdmin, models.RoleDev)).Post("/universities/{universityId}/courses", catalogController.CreateCourse)
		protected.With(authmiddleware.RequireRoles(models.RoleAdmin, models.RoleDev)).Post("/courses/{courseId}/subjects", catalogController.CreateSubject)

		protected.With(authmiddleware.RequireRoles(models.RoleApprover, models.RoleAdmin, models.RoleDev)).Get("/exams/pending", examController.ListPending)
		protected.With(authmiddleware.RequireRoles(models.RoleUser, models.RoleApprover, models.RoleAdmin, models.RoleDev)).Get("/users/me/exams/pending", examController.ListMyPending)
		protected.With(authmiddleware.RequireRoles(models.RoleUser, models.RoleApprover, models.RoleAdmin, models.RoleDev)).Post("/exams", examController.Upload)
		protected.With(authmiddleware.RequireRoles(models.RoleApprover, models.RoleAdmin, models.RoleDev)).Patch("/exams/{examId}/status", examController.Review)

		protected.With(authmiddleware.RequireRoles(models.RoleUser, models.RoleApprover, models.RoleAdmin, models.RoleDev)).Get("/users/me", userController.Me)
		protected.With(authmiddleware.RequireRoles(models.RoleAdmin, models.RoleDev)).Get("/users", userController.ListAll)
		protected.With(authmiddleware.RequireRoles(models.RoleAdmin, models.RoleDev)).Patch("/users/{id}/role", userController.UpdateRole)
	})

	return router, func() {}, nil
}

func health(w http.ResponseWriter, r *http.Request) {
	httpx.RespondJSON(w, http.StatusOK, map[string]string{"status": "UP"})
}
