package repository

import (
	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"gorm.io/gorm"
)

type BOMRepository interface {
	CreateBOM(bom *entity.Bom) (*entity.Bom, error)
	GetLastBomId() (string, error)
	CheckProductExists(productId string) (bool, error)
}

type bomRepository struct {
	db *gorm.DB
}

func NewBOMRepository(db *gorm.DB) BOMRepository {
	return &bomRepository{db: db}
}

func (r *bomRepository) CreateBOM(bom *entity.Bom) (*entity.Bom, error) {
	if err := r.db.Create(bom).Error; err != nil {
		return nil, err
	}
	return bom, nil
}

func (r *bomRepository) GetLastBomId() (string, error) {
	var lastBom entity.Bom
	err := r.db.Order("id_bom desc").First(&lastBom).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	return lastBom.BomId, nil
}

func (r *bomRepository) CheckProductExists(productId string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Products{}).Where("id_product = ?", productId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

