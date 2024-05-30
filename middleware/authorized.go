package middleware

import (
	"go-crud/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, exists := c.Get("Auth")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		u := auth.(models.Auth)

		for _, role := range u.Roles {
			for _, r := range roles {
				if role.Name == r {
					c.Next()
					return
				}
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
		c.Abort()
	}
}
