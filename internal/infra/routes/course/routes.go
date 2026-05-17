package courseroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/course"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/middleware"
)

func RegisterPublic(router *gin.Engine, controller *coursecontroller.CourseController) {
	router.GET("/universities/:universityId/courses", controller.List)
}
func RegisterProtected(router *gin.RouterGroup, controller *coursecontroller.CourseController) {
	router.POST("/universities/:universityId/courses", middleware.RequireRoles(user.RoleAdmin, user.RoleDev), controller.Create)
}
