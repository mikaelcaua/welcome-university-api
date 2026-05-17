package examroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/exam"
	"github.com/mikaelcaua/welcome-university-api/internal/domain/entities/user"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/middleware"
)

func RegisterPublic(router *gin.Engine, controller *examcontroller.ExamController) {
	router.GET("/subjects/:subjectId/exams", controller.ListBySubject)
	router.GET("/exams", controller.ListAll)
}
func RegisterProtected(router *gin.RouterGroup, controller *examcontroller.ExamController) {
	router.GET("/exams/pending", middleware.RequireRoles(user.RoleApprover, user.RoleAdmin, user.RoleDev), controller.ListPending)
	router.GET("/users/me/exams/pending", middleware.RequireRoles(user.RoleUser, user.RoleApprover, user.RoleAdmin, user.RoleDev), controller.ListMyPending)
	router.POST("/exams", middleware.RequireRoles(user.RoleUser, user.RoleApprover, user.RoleAdmin, user.RoleDev), controller.Upload)
	router.PATCH("/exams/:examId/status", middleware.RequireRoles(user.RoleApprover, user.RoleAdmin, user.RoleDev), controller.Review)
}
