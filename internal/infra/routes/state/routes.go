package stateroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/state"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/middleware"
)

func RegisterPublic(router *gin.Engine, controller *statecontroller.StateController) {
	router.GET("/states", controller.List)
	router.GET("/states/:code", controller.GetByCode)
}
func RegisterProtected(router *gin.RouterGroup, controller *statecontroller.StateController) {
	router.POST("/states", middleware.RequireRoles(user.RoleAdmin, user.RoleDev), controller.Create)
}
