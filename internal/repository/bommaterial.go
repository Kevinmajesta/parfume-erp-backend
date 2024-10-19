package repository

import (
	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"gorm.io/gorm"
)

type BOMMaterialRepository interface {
	CreateMaterial(material *entity.BomMaterial) (*entity.BomMaterial, error)
	GetLastMaterialId() (string, error)
	GetMaterialsByBomId(bomId string) ([]entity.BomMaterial, error)
	CheckMaterialExists(materialId string) (bool, error)
}

type bomMaterialRepository struct {
	db *gorm.DB
}

func NewBOMMaterialRepository(db *gorm.DB) BOMMaterialRepository {
	return &bomMaterialRepository{db: db}
}

func (r *bomMaterialRepository) CreateMaterial(material *entity.BomMaterial) (*entity.BomMaterial, error) {
	if err := r.db.Create(material).Error; err != nil {
		return nil, err
	}
	return material, nil
}

func (r *bomMaterialRepository) GetLastMaterialId() (string, error) {
	var lastMaterial entity.BomMaterial
	err := r.db.Order("id_bommaterial DESC").First(&lastMaterial).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	return lastMaterial.IdBomMaterial, nil
}

func (r *bomMaterialRepository) GetMaterialsByBomId(bomId string) ([]entity.BomMaterial, error) {
	var materials []entity.BomMaterial
	if err := r.db.Where("id_bom = ?", bomId).Find(&materials).Error; err != nil {
		return nil, err
	}
	return materials, nil
}

func (r *bomMaterialRepository) CheckMaterialExists(materialId string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Materials{}).Where("id_material = ?", materialId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
