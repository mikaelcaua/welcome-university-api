package subjectroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/subject"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/middleware"
)

func RegisterPublic(router *gin.Engine, controller *subjectcontroller.SubjectController) {
	router.GET("/courses/:courseId/subjects", controller.List)
}
func RegisterProtected(router *gin.RouterGroup, controller *subjectcontroller.SubjectController) {
	router.POST("/courses/:courseId/subjects", middleware.RequireRoles(user.RoleAdmin, user.RoleDev), controller.Create)
}
