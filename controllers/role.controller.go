package controllers

import (
	"go-crud/helpers"
	"go-crud/models"
	"go-crud/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RoleResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"title"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type RoleController struct {
	repo repositories.IRoleRepository
}

func NewRoleController(repo repositories.IRoleRepository) *RoleController {
	return &RoleController{repo}
}

func (ctrl *RoleController) CreateRole(c *gin.Context) {
	var roleData struct {
		Name string
	}

	c.Bind(&roleData)

	role := models.Role{
		Name: roleData.Name,
		Slug: helpers.FormatTextToSlug(roleData.Name),
	}

	if err := ctrl.repo.Create(&role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		Slug:      role.Slug,
		CreatedAt: role.CreatedAt,
	})
}

func (ctrl *RoleController) GetAllRoles(c *gin.Context) {
	roles, err := ctrl.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (ctrl *RoleController) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")

	// Convert the string ID to an unsigned integer
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	roleID := uint(id)

	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role.ID = roleID

	if err := ctrl.repo.Update(&role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

func (ctrl *RoleController) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")

	// Convert the string ID to an unsigned integer
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	roleID := uint(id)

	if err := ctrl.repo.Delete(roleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}
