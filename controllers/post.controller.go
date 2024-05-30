package controllers

import (
	"fmt"
	"go-crud/helpers"
	"go-crud/models"
	"go-crud/repositories"
	"go-crud/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	repo repositories.IPostRepository
}

func NewPostController(repo repositories.IPostRepository) *PostController {
	return &PostController{repo: repo}
}

func (ctrl *PostController) PostCreate(c *gin.Context) {
	idStr, ok := c.Get("UserID")
	if !ok || idStr == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID not found"})
		return
	}

	idStrVal, ok := idStr.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID is not a valid string"})
		return
	}

	id, err := strconv.ParseUint(idStrVal, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}
	userID := uint(id)

	// Define a struct to bind the incoming JSON
	var postData struct {
		Title string `json:"title" binding:"required"`
		Body  string `json:"body" binding:"required"`
	}

	// Bind the JSON to the postData struct and handle errors
	if err := c.ShouldBindJSON(&postData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the post model
	post := models.Post{
		Title:    postData.Title,
		Body:     postData.Body,
		Slug:     helpers.FormatTextToSlug(postData.Title),
		AuthorID: userID,
	}

	// Save the post to the repository and handle errors
	if err := ctrl.repo.Create(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%v has been created!", post.Title),
	})
}

func (ctrl *PostController) PostIndex(c *gin.Context) {
	limit := 10
	page := 1

	if limitQuery, exists := c.GetQuery("limit"); exists {
		if parsedLimit, err := strconv.Atoi(limitQuery); err == nil {
			limit = parsedLimit
		}
	}

	if pageQuery, exists := c.GetQuery("page"); exists {
		if parsedPage, err := strconv.Atoi(pageQuery); err == nil {
			page = parsedPage
		}
	}

	offset := (page - 1) * limit

	posts, total, err := ctrl.repo.FindAll(limit, offset)
	totalPages := (total + limit - 1) / limit

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.PaginatedResponse{
		Data:        posts,
		Page:        page,
		PageSize:    limit,
		Total:       total,
		Total_Pages: totalPages,
	})
}

func (ctrl *PostController) PostShow(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": fmt.Sprintf("UserID %v is invalid", id)})
		return
	}
	userID := uint(id)

	post, err := ctrl.repo.FindByID(userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "post not found!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, types.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Body:      post.Body,
		Slug:      post.Slug,
		CreatedAt: post.CreatedAt,
		AuthID:    post.AuthorID,
		Auth: types.AuthResponse{
			ID:    post.AuthorID,
			Email: post.Author.Email,
		},
	})
}

func (ctrl *PostController) PostUpdate(c *gin.Context) {
	idStr := c.Param("id")

	// Convert the string ID to an unsigned integer
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID := uint(id)

	post, err := ctrl.repo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "post not found", "error": err.Error()})
		return
	}

	var updatedPost models.Post
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if post.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You are not allowed to edit this post!"})
		return
	}

	post.Title = updatedPost.Title
	post.Slug = helpers.FormatTextToSlug(updatedPost.Title)
	post.Body = updatedPost.Body
	post.AuthorID = userID

	if err := ctrl.repo.Update(post); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "update failed!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Post %s has been updated!", post.Title),
	})
}

func (ctrl *PostController) PostDelete(c *gin.Context) {
	idStr := c.Param("id")

	// Convert the string ID to an unsigned integer
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID := uint(id)
	post, err := ctrl.repo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Item not found please try again!",
			"error":   err.Error(),
		})
		return
	}

	if err := ctrl.repo.Delete(post); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Something went wrong!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("Post %s has been deleted!", post.Title),
	})
}
