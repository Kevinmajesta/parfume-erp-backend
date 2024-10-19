package repository

import (
	"errors"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"gorm.io/gorm"
)

type BOMMaterialRepository interface {
	CreateMaterial(material *entity.BomMaterial) (*entity.BomMaterial, error)
	GetLastMaterialId() (string, error)
	GetMaterialsByBomId(bomId string) ([]entity.BomMaterial, error)
	CheckMaterialExists(materialId string) (bool, error)
	DeleteMaterialsByBomId(bomId string) error
	FindBOMByMaterialIDAndBOMID(materialId string, bomId string) (*entity.BomMaterial, error)
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

func (r *bomMaterialRepository) FindBOMByMaterialIDAndBOMID(materialId string, bomId string) (*entity.BomMaterial, error) {
	var bom entity.BomMaterial
	if err := r.db.Where("id_material = ? AND id_bom = ?", materialId, bomId).First(&bom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Tidak ditemukan duplikasi
		}
		return nil, err // Error lain saat pengecekan
	}
	return &bom, nil // Jika ditemukan, artinya ada duplikasi
}

func (r *bomMaterialRepository) DeleteMaterialsByBomId(bomId string) error {
	// Unscoped delete (hard delete) for bom materials related to bomId
	if err := r.db.Unscoped().Where("id_bom = ?", bomId).Delete(&entity.BomMaterial{}).Error; err != nil {
		return err
	}
	return nil
}
