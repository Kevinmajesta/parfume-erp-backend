package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/cache"
	"gorm.io/gorm"
)

type MaterialRepository interface {
	CreateMaterial(material *entity.Materials) (*entity.Materials, error)
	GetLastMaterial() (string, error)
	CheckMaterialExists(materialId string) (bool, error)
	UpdateMaterial(material *entity.Materials) (*entity.Materials, error)
	FindMaterialByID(materialId string) (*entity.Materials, error)
	DeleteMaterial(material *entity.Materials) (bool, error)
	FindAllMaterial(page int) ([]entity.Materials, error)
	SearchByName(name string) ([]entity.Materials, error)
}

type materialRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewMaterialRepository(db *gorm.DB, cacheable cache.Cacheable) *materialRepository {
	return &materialRepository{db: db, cacheable: cacheable}
}

func (r *materialRepository) CreateMaterial(material *entity.Materials) (*entity.Materials, error) {
	if err := r.db.Create(&material).Error; err != nil {
		return material, err
	}
	r.cacheable.Delete("FindAllMaterials_page_1")
	r.cacheable.Delete("FindAllMaterials_page_2")
	return material, nil
}

func (r *materialRepository) GetLastMaterial() (string, error) {
	var lastMaterial entity.Materials
	err := r.db.Order("id_material DESC").First(&lastMaterial).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Jika tidak ada produk terakhir, return ID default
		return "MTR-00000", nil
	} else if err != nil {
		return "", err
	}

	// Jika produk ditemukan, return ProdukId dari produk terakhir
	return lastMaterial.MaterialId, nil
}

func (r *materialRepository) CheckMaterialExists(materialId string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Materials{}).Where("id_material = ?", materialId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *materialRepository) UpdateMaterial(material *entity.Materials) (*entity.Materials, error) {
	fields := make(map[string]interface{})

	if material.Materialname != "" {
		fields["materialname"] = material.Materialname
	}
	if material.Materialcategory != "" {
		fields["materialcategory"] = material.Materialcategory
	}
	if material.Sellprice != "" {
		fields["sellprice"] = material.Sellprice
	}
	if material.Makeprice != "" {
		fields["makeprice"] = material.Makeprice
	}
	if material.Unit != "" {
		fields["unit"] = material.Unit
	}
	if material.Description != "" {
		fields["description"] = material.Description
	}
	if material.Image != "" {
		fields["image"] = material.Image
	}

	if err := r.db.Model(material).Where("id_material = ?", material.MaterialId).Updates(fields).Error; err != nil {
		return material, err
	}
	r.cacheable.Delete("FindAllMaterials_page_1")

	return material, nil
}

func (r *materialRepository) FindMaterialByID(materialId string) (*entity.Materials, error) {
	material := new(entity.Materials)
	if err := r.db.Where("id_material = ?", materialId).First(material).Error; err != nil {
		log.Printf("Error finding material by ID: %v", err)
		return nil, err // Pastikan mengembalikan nil, err
	}
	log.Printf("material found: %v", material)
	return material, nil
}

func (r *materialRepository) DeleteMaterial(material *entity.Materials) (bool, error) {
	log.Printf("Deleting material: %v", material)
	// Ensure hard delete by using Unscoped()
	if err := r.db.Unscoped().Delete(material).Error; err != nil {
		log.Printf("Error deleting material: %v", err)
		return false, err
	}
	log.Println("material deleted successfully")
	r.cacheable.Delete("FindAllMaterials_page_1")
	return true, nil
}

func (r *materialRepository) FindAllMaterial(page int) ([]entity.Materials, error) {
	var Materials []entity.Materials
	key := fmt.Sprintf("FindAllMaterials_page_%d", page)
	const pageSize = 100

	data, _ := r.cacheable.Get(key)
	if data == "" {
		offset := (page - 1) * pageSize
		if err := r.db.Limit(pageSize).Offset(offset).Find(&Materials).Error; err != nil {
			return Materials, err
		}
		marshalledMaterials, _ := json.Marshal(Materials)
		err := r.cacheable.Set(key, marshalledMaterials, 5*time.Minute)
		if err != nil {
			return Materials, err
		}
	} else {
		err := json.Unmarshal([]byte(data), &Materials)
		if err != nil {
			return Materials, err
		}
	}
	return Materials, nil
}

func (r *materialRepository) SearchByName(name string) ([]entity.Materials, error) {
	var Materials []entity.Materials
	query := r.db
	// Gunakan fungsi LOWER untuk mengabaikan perbedaan huruf besar dan kecil
	if err := query.Where("LOWER(materialname) LIKE LOWER(?)", "%"+name+"%").Find(&Materials).Error; err != nil {
		return nil, err
	}
	return Materials, nil
}
