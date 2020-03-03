package repository

import (
	"api/models"
)

type UserRepository interface {
	Save(models.User) (models.User, error)
	FindAll() ([]models.User, error)
	FindByID(uint32) (models.User, error)
	// update(uint32, models.User) (int64, error) /// why int64?
	// Delete(uint32) (int64, error)
}
