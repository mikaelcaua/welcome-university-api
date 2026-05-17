package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/auth"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/course"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/state"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/university"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/user"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/httpx"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/auth"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/course"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/state"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/university"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/routes/user"
)

type Dependencies struct {
	AuthController       *authcontroller.AuthController
	UserController       *usercontroller.UserController
	StateController      *statecontroller.StateController
	UniversityController *universitycontroller.UniversityController
	CourseController     *coursecontroller.CourseController
	SubjectController    *subjectcontroller.SubjectController
	ExamController       *examcontroller.ExamController
	Authentication       gin.HandlerFunc
}

func Register(router *gin.Engine, dependencies Dependencies) {
	router.GET("/actuator/health/readiness", health)
	router.GET("/actuator/health/liveness", health)
	router.GET("/actuator/info", func(ctx *gin.Context) {
		httpx.RespondJSON(ctx, http.StatusOK, gin.H{"app": "welcome-university-api"})
	})

	authroutes.Register(router, dependencies.AuthController)
	stateroutes.RegisterPublic(router, dependencies.StateController)
	universityroutes.RegisterPublic(router, dependencies.UniversityController)
	courseroutes.RegisterPublic(router, dependencies.CourseController)
	subjectroutes.RegisterPublic(router, dependencies.SubjectController)
	examroutes.RegisterPublic(router, dependencies.ExamController)

	protected := router.Group("", dependencies.Authentication)
	stateroutes.RegisterProtected(protected, dependencies.StateController)
	universityroutes.RegisterProtected(protected, dependencies.UniversityController)
	courseroutes.RegisterProtected(protected, dependencies.CourseController)
	subjectroutes.RegisterProtected(protected, dependencies.SubjectController)
	examroutes.RegisterProtected(protected, dependencies.ExamController)
	userroutes.RegisterProtected(protected, dependencies.UserController)
}

func health(ctx *gin.Context) {
	httpx.RespondJSON(ctx, http.StatusOK, gin.H{"status": "UP"})
}
