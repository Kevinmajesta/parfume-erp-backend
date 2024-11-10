package repository

import (
	"errors"
	"log"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/cache"
	"gorm.io/gorm"
)

type VendorRepository interface {
	CreateVendor(vendor *entity.Vendors) (*entity.Vendors, error)
	GetLastVendor() (string, error)
	CheckVendorExists(vendorId string) (bool, error)
	UpdateVendor(vendor *entity.Vendors) (*entity.Vendors, error)
	FindVendorByID(vendorId string) (*entity.Vendors, error)
	DeleteVendor(vendor *entity.Vendors) (bool, error)
}

type vendorRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewVendorRepository(db *gorm.DB, cacheable cache.Cacheable) *vendorRepository {
	return &vendorRepository{db: db, cacheable: cacheable}
}

func (r *vendorRepository) CreateVendor(vendor *entity.Vendors) (*entity.Vendors, error) {
	if err := r.db.Create(&vendor).Error; err != nil {
		return vendor, err
	}
	r.cacheable.Delete("FindAllVendors_page_1")
	r.cacheable.Delete("FindAllVendors_page_2")
	return vendor, nil
}

func (r *vendorRepository) GetLastVendor() (string, error) {
	var lastVendor entity.Vendors
	err := r.db.Order("id_vendor DESC").First(&lastVendor).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "VDR-00000", nil
	} else if err != nil {
		return "", err
	}

	return lastVendor.VendorId, nil
}

func (r *vendorRepository) CheckVendorExists(vendorId string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Vendors{}).Where("id_vendor = ?", vendorId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *vendorRepository) UpdateVendor(vendor *entity.Vendors) (*entity.Vendors, error) {
	fields := make(map[string]interface{})

	if vendor.Vendorname != "" {
		fields["vendorname"] = vendor.Vendorname
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
	if vendor.Website != "" {
		fields["website"] = vendor.Website
	}

	if err := r.db.Model(vendor).Where("id_vendor = ?", vendor.VendorId).Updates(fields).Error; err != nil {
		return vendor, err
	}
	r.cacheable.Delete("FindAllPVendors_page_1")

	return vendor, nil
}

func (r *vendorRepository) FindVendorByID(vendorId string) (*entity.Vendors, error) {
	vendors := new(entity.Vendors)
	if err := r.db.Where("id_vendor = ?", vendorId).First(vendors).Error; err != nil {
		return nil, err
	}
	return vendors, nil
}

func (r *vendorRepository) DeleteVendor(vendor *entity.Vendors) (bool, error) {
	if err := r.db.Unscoped().Delete(vendor).Error; err != nil {
		log.Printf("Error deleting vendor: %v", err)
		return false, err
	}
	r.cacheable.Delete("FindAllVendors_page_1")
	return true, nil
}
