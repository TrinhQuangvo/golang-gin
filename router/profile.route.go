package router

import (
	"go-crud/controllers"
	"go-crud/initilizers"
	"go-crud/middleware"
	"go-crud/repositories"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(r *gin.Engine) {
	profileRepo := repositories.NewProfileRepository(initilizers.DB)
	profileController := controllers.NewProfileController(profileRepo)
	profileGroup := r.Group("/profile")
	{
		profileGroup.GET("/", middleware.AuthRequired, profileController.GetProfile)
		profileGroup.PATCH("/:id", middleware.AuthRequired, profileController.UpdateProfile)
		profileGroup.POST("/create-new-profile", middleware.AuthRequired, profileController.CreateNewProfile)
	}
}
