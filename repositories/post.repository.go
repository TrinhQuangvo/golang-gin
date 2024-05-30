package repositories

import (
	"go-crud/models"

	"gorm.io/gorm"
)

type IPostRepository interface {
	Create(post *models.Post) error
	FindAll(limit int, take int) ([]models.Post, int, error)
	FindByID(id uint) (*models.Post, error)
	Update(post *models.Post) error
	Delete(post *models.Post) error
}
type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) IPostRepository {
	return &postRepository{db}
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindAll(limit, offset int) ([]models.Post, int, error) {
	var posts []models.Post
	var total int64

	r.db.Model(&models.Post{}).Count(&total)

	result := r.db.Model(&models.Post{}).Preload("Author").Find(&posts).Error

	return posts, int(total), result
}

func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post

	err := r.db.Preload("Author").Where("id = ?", id).First(&post, id).Error

	return &post, err
}

func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(post *models.Post) error {
	return r.db.Delete(post).Error
}
