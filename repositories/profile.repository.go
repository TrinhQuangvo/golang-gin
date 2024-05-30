package repositories

import (
	"go-crud/models"

	"gorm.io/gorm"
)

type IProfileRepository interface {
	Create(profile *models.Profile) error
	ShowProfile(id uint) (*models.Profile, error)
	Update(profile *models.Profile) error
}
type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) IProfileRepository {
	return &profileRepository{db: db}
}

func (repo *profileRepository) Create(profile *models.Profile) error {
	return repo.db.Create(&profile).Error
}

func (repo *profileRepository) ShowProfile(id uint) (*models.Profile, error) {
	var profile models.Profile
	err := repo.db.First(&profile, id).Error
	return &profile, err
}

func (repo *profileRepository) Update(profile *models.Profile) error {
	return repo.db.Save(profile).Error
}
