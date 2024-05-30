package router

import (
	"go-crud/controllers"
	"go-crud/initilizers"
	"go-crud/repositories"

	"github.com/gin-gonic/gin"
)

func RoleRoutes(r *gin.Engine) {
	roleRepo := repositories.NewRoleRepository(initilizers.DB)
	roleController := controllers.NewRoleController(roleRepo)

	roleGroup := r.Group("/roles")
	{
		roleGroup.POST("/", roleController.CreateRole)
		roleGroup.GET("/", roleController.GetAllRoles)
		roleGroup.PATCH("/:id", roleController.UpdateRole)
		roleGroup.DELETE("/:id", roleController.DeleteRole)
	}
}
