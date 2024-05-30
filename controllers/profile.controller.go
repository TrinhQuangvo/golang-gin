package controllers

import (
	"fmt"
	"go-crud/models"
	"go-crud/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	repo repositories.IProfileRepository
}

func NewProfileController(repo repositories.IProfileRepository) *ProfileController {
	return &ProfileController{repo: repo}
}

func (ctrl *ProfileController) GetProfile(c *gin.Context) {
	auth, _ := c.Get("Auth")
	user := auth.(models.Auth)

	if auth != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": error.Error,
		})
		return
	}

	profile, err := ctrl.repo.ShowProfile(user.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"profile": profile,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"profile": profile,
	})
}
func (ctrl *ProfileController) CreateNewProfile(c *gin.Context) {
	var profileData struct {
		FriendlyName string
		Status       bool
		Address      string
		PhoneNumber  string
		Avatar       string
	}
	c.Bind(&profileData)

	profile := models.Profile{
		FriendlyName: profileData.FriendlyName,
		Address:      profileData.Address,
		PhoneNumber:  profileData.PhoneNumber,
		Avartar:      profileData.Avatar,
		Status:       true,
	}

	if err := c.ShouldBindJSON(profileData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed!",
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"profile": profile,
	})
}

func (ctrl *ProfileController) UpdateProfile(c *gin.Context) {
	var profileData struct {
		FriendlyName string
		Address      string
		PhoneNumber  string
		Avatar       string
	}
	c.ShouldBindJSON(&profileData)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID := uint(id)
	profile, err := ctrl.repo.ShowProfile(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No profile found!",
		})
		return
	}

	if err := c.Bind(&profileData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Something went wrong!",
			"error":   err.Error(),
		})
		return
	}

	profile.Address = profileData.Address
	profile.FriendlyName = profileData.FriendlyName
	profile.PhoneNumber = profileData.PhoneNumber

	if err := ctrl.repo.Update(profile); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "update profile failed!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("update profile %s successfully!", profile.FriendlyName),
	})
}
