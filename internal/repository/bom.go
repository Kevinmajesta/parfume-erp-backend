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

type BOMRepository interface {
	CreateBOM(bom *entity.Bom) (*entity.Bom, error)
	GetLastBomId() (string, error)
	CheckProductExists(productId string) (bool, error)
	FindAllBom(page int) ([]entity.Bom, error)
	DeleteBom(bomId string) (bool, error)
	UpdateBOM(bom *entity.Bom) (*entity.Bom, error)
	FindBOMByProductIDAndBOMID(productId string, bomId string) (*entity.Bom, error)
	FindBOMByID(bomId string) (*entity.Bom, error)
	GetProductDetails(productId string) (*entity.Products, error)
}

type bomRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewBOMRepository(db *gorm.DB, cacheable cache.Cacheable) BOMRepository {
	return &bomRepository{db: db, cacheable: cacheable}
}

func (r *bomRepository) CreateBOM(bom *entity.Bom) (*entity.Bom, error) {
	if err := r.db.Create(bom).Error; err != nil {
		return nil, err
	}
	r.cacheable.Delete("FindAllBoms_page_1")
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

func (r *bomRepository) FindAllBom(page int) ([]entity.Bom, error) {
	var Bom []entity.Bom
	key := fmt.Sprintf("FindAllBoms_page_%d", page)
	const pageSize = 100

	data, _ := r.cacheable.Get(key)
	if data == "" {
		offset := (page - 1) * pageSize
		if err := r.db.Limit(pageSize).Offset(offset).Find(&Bom).Error; err != nil {
			return Bom, err
		}
		marshalledBoms, _ := json.Marshal(Bom)
		err := r.cacheable.Set(key, marshalledBoms, 5*time.Minute)
		if err != nil {
			return Bom, err
		}
	} else {
		err := json.Unmarshal([]byte(data), &Bom)
		if err != nil {
			return Bom, err
		}
	}
	return Bom, nil
}

func (r *bomRepository) DeleteBom(bomId string) (bool, error) {

	if err := r.db.Unscoped().Where("id_bom = ?", bomId).Delete(&entity.BomMaterial{}).Error; err != nil {
		log.Printf("Error deleting materials for bom ID %s: %v", bomId, err)
		return false, err
	}

	if err := r.db.Unscoped().Where("id_bom = ?", bomId).Delete(&entity.Bom{}).Error; err != nil {
		log.Printf("Error deleting bom with ID %s: %v", bomId, err)
		return false, err
	}

	log.Printf("Successfully deleted bom with ID %s and its related materials", bomId)
	r.cacheable.Delete("FindAllBoms_page_1")
	return true, nil
}

func (r *bomRepository) UpdateBOM(bom *entity.Bom) (*entity.Bom, error) {
	// Update BOM table
	if err := r.db.Model(&entity.Bom{}).Where("id_bom = ?", bom.BomId).Updates(bom).Error; err != nil {
		return nil, err
	}
	// Clear cache after update
	r.cacheable.Delete("FindAllBoms_page_1")
	return bom, nil
}

func (r *bomRepository) FindBOMByProductIDAndBOMID(productId string, bomId string) (*entity.Bom, error) {
	var bom entity.Bom
	if err := r.db.Where("id_product = ? AND id_bom = ?", productId, bomId).First(&bom).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Tidak ditemukan duplikasi
		}
		return nil, err // Error lain saat pengecekan
	}
	return &bom, nil // Jika ditemukan, artinya ada duplikasi
}

func (r *bomRepository) FindBOMByID(bomId string) (*entity.Bom, error) {
    var bom entity.Bom
    if err := r.db.Preload("Materials").Where("id_bom = ?", bomId).First(&bom).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil 
        }
        return nil, err
    }
    return &bom, nil 
}

func (r *bomRepository) GetProductDetails(productId string) (*entity.Products, error) {
    var product entity.Products
    if err := r.db.Table("products").Where("id_product = ?", productId).First(&product).Error; err != nil {
        return nil, err
    }
    return &product, nil
}


