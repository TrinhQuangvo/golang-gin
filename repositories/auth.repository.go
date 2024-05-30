package repositories

import (
	"go-crud/models"

	"gorm.io/gorm"
)

type IAuthRepository interface {
	FindByEmail(email string) (*models.Auth, error)
	Create(auth *models.Auth) error
	FindByRefreshToken(token string) (*models.Auth, error)
	Save(auth *models.Auth) error
	ChangePassword(id uint, newPassword string) error
	FindRoleByID(id uint) (*models.Role, error)
	FindAllUsers(limit int, offset int) ([]models.Auth, int, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindByEmail(email string) (*models.Auth, error) {
	var auth models.Auth
	if err := r.db.Where("email = ?", email).First(&auth).Select("auths.ID, auths.email, auths.body,").Joins("left join profile on profile_id = auths.profile_id").Error; err != nil {
		return nil, err
	}
	return &auth, nil
}

func (r *authRepository) FindAllUsers(limit, offset int) ([]models.Auth, int, error) {
	var auths []models.Auth
	var total int64

	r.db.Model(&models.Auth{}).Count(&total)

	err := r.db.Preload("Roles").Find(&auths).Error

	return auths, int(total), err
}

func (r *authRepository) FindRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *authRepository) Create(auth *models.Auth) error {
	return r.db.Create(auth).Error
}

func (r *authRepository) FindByRefreshToken(token string) (*models.Auth, error) {
	var auth models.Auth
	err := r.db.Where("refresh_token = ?", token).First(&auth).Error
	return &auth, err
}

func (r *authRepository) Save(auth *models.Auth) error {
	return r.db.Save(auth).Error
}

func (r *authRepository) ChangePassword(id uint, newPassword string) error {
	var auth models.Auth
	if err := r.db.First(&auth, id).Error; err != nil {
		return err
	}
	auth.Password = newPassword
	return r.db.Save(&auth).Error
}
