package router

import (
	"go-crud/controllers"
	"go-crud/initilizers"
	"go-crud/middleware"
	"go-crud/repositories"

	"github.com/gin-gonic/gin"
)

func PostRoutes(r *gin.Engine) {
	postRepo := repositories.NewPostRepository(initilizers.DB)
	postController := controllers.NewPostController(postRepo)
	postGroup := r.Group("/post")
	{
		postGroup.POST("/", middleware.AuthRequired, postController.PostCreate)
		postGroup.GET("/", postController.PostIndex)
		postGroup.GET("/:id", postController.PostShow)
		postGroup.PATCH("/:id", middleware.AuthRequired, postController.PostUpdate)
		postGroup.DELETE("/:id", middleware.AuthRequired, postController.PostDelete)
	}
}
