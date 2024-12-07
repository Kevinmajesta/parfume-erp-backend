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

type CostumerRepository interface {
	CreateCostumer(vendor *entity.Costumers) (*entity.Costumers, error)
	GetLastCostumer() (string, error)
	CheckCostumerExists(vendorId string) (bool, error)
	UpdateCostumer(vendor *entity.Costumers) (*entity.Costumers, error)
	DeleteCostumer(vendor *entity.Costumers) (bool, error)
	FindCostumerByID(vendorId string) (*entity.Costumers, error)
	FindAllCostumer(page int) ([]entity.Costumers, error)
}

type costumerRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewCostumerRepository(db *gorm.DB, cacheable cache.Cacheable) *costumerRepository {
	return &costumerRepository{db: db, cacheable: cacheable}
}

func (r *costumerRepository) CreateCostumer(vendor *entity.Costumers) (*entity.Costumers, error) {
	if err := r.db.Create(&vendor).Error; err != nil {
		return vendor, err
	}
	r.cacheable.Delete("FindAllCostumers_page_1")
	return vendor, nil
}

func (r *costumerRepository) GetLastCostumer() (string, error) {
	var lastVendor entity.Costumers
	err := r.db.Order("id_costumer DESC").First(&lastVendor).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "CRS-00000", nil
	} else if err != nil {
		return "", err
	}

	return lastVendor.CostumerId, nil
}

func (r *costumerRepository) CheckCostumerExists(vendorId string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Costumers{}).Where("id_costumer = ?", vendorId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *costumerRepository) UpdateCostumer(vendor *entity.Costumers) (*entity.Costumers, error) {
	fields := make(map[string]interface{})

	if vendor.Costumername != "" {
		fields["costumername"] = vendor.Costumername
	}
	if vendor.Addressone != "" {
		fields["addressone"] = vendor.Addressone
	}
	if vendor.Addresstwo != "" {
		fields["addresstwo"] = vendor.Addresstwo
	}
	if vendor.Phone != "" {
		fields["phone"] = vendor.Phone
	}
	if vendor.Email != "" {
		fields["email"] = vendor.Email
	}

	if err := r.db.Model(vendor).Where("id_costumer = ?", vendor.CostumerId).Updates(fields).Error; err != nil {
		return vendor, err
	}
	r.cacheable.Delete("FindAllCostumers_page_1")

	return vendor, nil
}

func (r *costumerRepository) FindCostumerByID(vendorId string) (*entity.Costumers, error) {
	vendors := new(entity.Costumers)
	if err := r.db.Where("id_costumer = ?", vendorId).First(vendors).Error; err != nil {
		return nil, err
	}
	return vendors, nil
}

func (r *costumerRepository) DeleteCostumer(vendor *entity.Costumers) (bool, error) {
	if err := r.db.Unscoped().Delete(vendor).Error; err != nil {
		log.Printf("Error deleting costumer: %v", err)
		return false, err
	}
	r.cacheable.Delete("FindAllCostumers_page_1")
	return true, nil
}

func (r *costumerRepository) FindAllCostumer(page int) ([]entity.Costumers, error) {
	var Vendors []entity.Costumers
	key := fmt.Sprintf("FindAllCostumers_page_%d", page)
	const pageSize = 100

	data, _ := r.cacheable.Get(key)
	if data == "" {
		offset := (page - 1) * pageSize
		if err := r.db.Limit(pageSize).Offset(offset).Find(&Vendors).Error; err != nil {
			return Vendors, err
		}
		marshalledVendors, _ := json.Marshal(Vendors)
		err := r.cacheable.Set(key, marshalledVendors, 5*time.Minute)
		if err != nil {
			return Vendors, err
		}
	} else {
		err := json.Unmarshal([]byte(data), &Vendors)
		if err != nil {
			return Vendors, err
		}
	}
	return Vendors, nil
}
