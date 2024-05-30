package router

import (
	"go-crud/controllers"
	"go-crud/initilizers"
	"go-crud/middleware"
	"go-crud/repositories"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	authRepo := repositories.NewAuthRepository(initilizers.DB)
	authController := controllers.NewAuthController(authRepo)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/sign-in", authController.SignIn)
		authGroup.GET("/", middleware.AuthRequired, authController.GetAllUsers)
		authGroup.POST("/sign-up", authController.SignUp)
		authGroup.GET("/validate", middleware.AuthRequired, authController.Validate)
		authGroup.POST("/log-out", middleware.AuthRequired, authController.Logout)
		authGroup.POST("/refresh", authController.RefreshToken)
		authGroup.PATCH("/change-password", middleware.AuthRequired, authController.ChangePassword)
	}
}
