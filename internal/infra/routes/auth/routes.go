package authroutes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/controllers/auth"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/middleware"
)

func Register(router *gin.Engine, controller *authcontroller.AuthController) {
	authRoutes := router.Group("/auth", middleware.NewRateLimiter(20, time.Minute).Middleware())
	authRoutes.POST("/register", controller.Register)
	authRoutes.POST("/login", controller.Login)
	authRoutes.POST("/refresh", controller.Refresh)
}
