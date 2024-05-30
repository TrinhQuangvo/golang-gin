package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	RoleRoutes(r)
	ProfileRoutes(r)
	AuthRoutes(r)
	PostRoutes(r)
	return r
}
