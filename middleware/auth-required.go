package middleware

import (
	"fmt"
	"go-crud/initilizers"
	"go-crud/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var auth models.Auth
		initilizers.DB.Preload("Roles").First(&auth, claims["sub"])
		if auth.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("Auth", auth)
		c.Set("UserID", auth.ID)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
