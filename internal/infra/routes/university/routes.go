package universityroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/university"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/middleware"
)

func RegisterPublic(router *gin.Engine, controller *universitycontroller.UniversityController) {
	router.GET("/states/:stateId/universities", controller.List)
}
func RegisterProtected(router *gin.RouterGroup, controller *universitycontroller.UniversityController) {
	router.POST("/states/:stateId/universities", middleware.RequireRoles(user.RoleAdmin, user.RoleDev), controller.Create)
}
