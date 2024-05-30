package controllers

import (
	"fmt"
	"go-crud/helpers"
	"go-crud/models"
	"go-crud/repositories"
	"go-crud/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	repo repositories.IAuthRepository
}

func NewAuthController(repo repositories.IAuthRepository) *AuthController {
	return &AuthController{repo: repo}
}

func (ctrl *AuthController) SignUp(c *gin.Context) {
	var authInfo struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		RoleID   []uint `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&authInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to read data!", "error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(authInfo.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password!", "error": err.Error()})
		return
	}

	auth := models.Auth{
		Email:    authInfo.Email,
		Password: string(hash),
	}

	// Find and assign roles
	for _, roleID := range authInfo.RoleID {
		role, err := ctrl.repo.FindRoleByID(roleID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Role not found!", "error": err.Error()})
			return
		}
		auth.Roles = append(auth.Roles, *role)
	}

	if err := ctrl.repo.Create(&auth); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Create user failed!", "auth": auth})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("account %s has been created!", auth.Email)})
}

func (ctrl *AuthController) SignIn(c *gin.Context) {
	var auth struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to read data!", "error": err.Error()})
		return
	}

	authModel, err := ctrl.repo.FindByEmail(auth.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password!"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(authModel.Password), []byte(auth.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password!"})
		return
	}

	tokenString, err := helpers.GenerateAccessToken(authModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create token!", "error": err.Error()})
		return
	}

	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
	c.SetSameSite(http.SameSiteDefaultMode)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (ctrl *AuthController) Validate(c *gin.Context) {
	id, _ := c.Get("UserID")
	auth, _ := c.Get("Auth")
	c.JSON(http.StatusOK, gin.H{"id": id, "auth": auth})
}

func (ctrl *AuthController) GetAllUsers(c *gin.Context) {
	_, ok := c.Get("UserID")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"Error": "You have to login first"})
		return
	}
	limit := 10
	offset := 1

	users, total, err := ctrl.repo.FindAllUsers(limit, offset)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.PaginatedResponse{
		Data:        users,
		Total:       total,
		PageSize:    limit,
		Total_Pages: offset,
	})
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authModel, err := ctrl.repo.FindByRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token"})
		return
	}

	accessToken, err := helpers.GenerateAccessToken(authModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create access token", "error": err.Error()})
		return
	}

	refreshToken, err := helpers.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create refresh token", "error": err.Error()})
		return
	}

	authModel.RefreshToken = refreshToken
	if err := ctrl.repo.Save(authModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}

func (ctrl *AuthController) ChangePassword(c *gin.Context) {
	var req struct {
		newPassword string
		oldPassword string
	}

	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect password!",
		})
		return
	}

	auth, _ := c.Get("Auth")
	user := auth.(models.Auth)

	if err := bcrypt.CompareHashAndPassword([]byte(req.newPassword), []byte(req.oldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Old password is incorrect"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.newPassword), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash new password"})
	}

	if err := ctrl.repo.ChangePassword(user.ID, string(hashedPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Change Password failed!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
