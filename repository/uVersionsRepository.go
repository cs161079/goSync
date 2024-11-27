package repository

import (
	models "github.com/cs161079/godbLib/Models"
	"gorm.io/gorm"
)

type UVersionRepository interface {
	UVersionList() ([]models.UVersions, error)
}

type uVersionRepository struct {
	DB *gorm.DB
}

func NewUVersionsRepository(connection *gorm.DB) UVersionRepository {
	return uVersionRepository{
		DB: connection,
	}
}

func (r uVersionRepository) UVersionList() ([]models.UVersions, error) {
	var result []models.UVersions
	if err := r.DB.Table("syncversions").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
