package user

import (
	"gorm.io/gorm"

	"github.com/saveblush/reraw-api/internal/repositories"
)

// repository interface
type Repository interface {
	FindByIDString(db *gorm.DB, field string, value string, i interface{}) error
}

type repository struct {
	repositories.Repository
}

func NewRepository() Repository {
	return &repository{
		repositories.NewRepository(),
	}
}
