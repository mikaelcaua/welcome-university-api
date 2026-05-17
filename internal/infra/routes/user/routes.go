package userroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/user"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/middleware"
)

func RegisterProtected(router *gin.RouterGroup, controller *usercontroller.UserController) {
	router.GET("/users/me", middleware.RequireRoles(user.RoleUser, user.RoleApprover, user.RoleAdmin, user.RoleDev), controller.Me)
	router.GET("/users", middleware.RequireRoles(user.RoleAdmin, user.RoleDev), controller.ListAll)
	router.PATCH("/users/:id/role", middleware.RequireRoles(user.RoleAdmin, user.RoleDev), controller.UpdateRole)
}
