package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/cache"
	"gorm.io/gorm"
)

type RfqRepository interface {
	GetLastRfq() (string, error)
	CreateRfq(mo *entity.Rfqs) (*entity.Rfqs, error)
	CheckProductExists(productId string) (bool, error)
	UpdateRfq(rfq *entity.Rfqs) (*entity.Rfqs, error)
	GetRfqById(rfqId string) (*entity.Rfqs, error)
	UpdateRfqStatus(rfq *entity.Rfqs) (*entity.Rfqs, error)
	FindAllRfq(page int) ([]entity.Rfqs, error)
	FindAllRfqBill(page int) ([]entity.Rfqs, error)
	GetVendorDetails(vendorId string) (*entity.Vendors, error)
	CheckEmailExistsByVendorId(vendorId string) (string, error)
}

type rfqRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewRfqRepository(db *gorm.DB, cacheable cache.Cacheable) *rfqRepository {
	return &rfqRepository{db: db, cacheable: cacheable}
}

func (r *rfqRepository) GetLastRfq() (string, error) {
	var lastRfq entity.Rfqs
	err := r.db.Order("id_rfq DESC").First(&lastRfq).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "R-00000", nil
	} else if err != nil {
		return "", err
	}

	return lastRfq.RfqId, nil
}

func (r *rfqRepository) CreateRfq(rfq *entity.Rfqs) (*entity.Rfqs, error) {
	if err := r.db.Create(&rfq).Error; err != nil {
		return rfq, err
	}
	r.cacheable.Delete("FindAllRfq_page_1")
	r.cacheable.Delete("FindAllRfqBill_page_1")
	return rfq, nil
}

func (r *rfqRepository) CheckProductExists(productId string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Vendors{}).Where("id_vendor = ?", productId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *rfqRepository) UpdateRfq(rfq *entity.Rfqs) (*entity.Rfqs, error) {
	if err := r.db.Model(&entity.Rfqs{}).Where("id_rfq = ?", rfq.RfqId).Updates(rfq).Error; err != nil {
		return nil, err
	}
	r.cacheable.Delete("FindAllBoms_page_1")
	r.cacheable.Delete("FindAllRfqBill_page_1")
	return rfq, nil
}

// func (r *rfqRepository) GetRfqById(rfqId string) (*entity.Rfqs, error) {
// 	var rfq entity.Rfqs
// 	if err := r.db.Where("id_rfq = ?", rfqId).First(&rfq).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	r.cacheable.Delete("FindAllBoms_page_1")
// 	return &rfq, nil
// }

func (r *rfqRepository) GetRfqById(rfqId string) (*entity.Rfqs, error) {
	var rfq entity.Rfqs
	if err := r.db.Where("id_rfq = ?", rfqId).
		Preload("Products"). // Memuat data dari tabel rfqs_products
		First(&rfq).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rfq, nil
}


func (r *rfqRepository) UpdateRfqStatus(rfq *entity.Rfqs) (*entity.Rfqs, error) {
	if err := r.db.Save(rfq).Error; err != nil {
		return nil, err
	}
	r.cacheable.Delete("FindAllRfq_page_1")
	r.cacheable.Delete("FindAllRfqBill_page_1")
	return rfq, nil
}

func (r *rfqRepository) FindAllRfq(page int) ([]entity.Rfqs, error) {
	var Rfq []entity.Rfqs
	key := fmt.Sprintf("FindAllRfq_page_%d", page)
	const pageSize = 100

	data, _ := r.cacheable.Get(key)
	if data == "" {
		offset := (page - 1) * pageSize
		if err := r.db.Where("status IN ?", []string{"RFQ", "Purchase Order"}).
			Limit(pageSize).
			Offset(offset).
			Find(&Rfq).Error; err != nil {
			return Rfq, err
		}
		marshalledRfqs, _ := json.Marshal(Rfq)
		err := r.cacheable.Set(key, marshalledRfqs, 5*time.Minute)
		if err != nil {
			return Rfq, err
		}
	} else {
		err := json.Unmarshal([]byte(data), &Rfq)
		if err != nil {
			return Rfq, err
		}
	}
	r.cacheable.Delete("FindAllRfq_page_1")
	r.cacheable.Delete("FindAllRfqBill_page_1")
	return Rfq, nil
}

func (r *rfqRepository) FindAllRfqBill(page int) ([]entity.Rfqs, error) {
	var Rfq []entity.Rfqs
	key := fmt.Sprintf("FindAllRfqBill_page_%d", page)
	const pageSize = 100

	data, _ := r.cacheable.Get(key)
	if data == "" {
		offset := (page - 1) * pageSize
		if err := r.db.Where("status IN ?", []string{"Recived", "Billed", "Done"}).
			Limit(pageSize).
			Offset(offset).
			Find(&Rfq).Error; err != nil {
			return Rfq, err
		}
		marshalledRfqs, _ := json.Marshal(Rfq)
		err := r.cacheable.Set(key, marshalledRfqs, 5*time.Minute)
		if err != nil {
			return Rfq, err
		}
	} else {
		err := json.Unmarshal([]byte(data), &Rfq)
		if err != nil {
			return Rfq, err
		}
	}
	r.cacheable.Delete("FindAllRfqBill_page_1")
	return Rfq, nil
}

func (r *rfqRepository) GetVendorDetails(vendorId string) (*entity.Vendors, error) {
	var product entity.Vendors
	if err := r.db.Table("products").Where("id_vendor = ?", vendorId).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *rfqRepository) CheckEmailExistsByVendorId(vendorId string) (string, error) {
	var vendor entity.Vendors
	err := r.db.Table("vendors").Where("id_vendor = ?", vendorId).Select("email").Scan(&vendor).Error
	if err != nil {
		return "", err
	}

	if vendor.Email == "" {
		return "", fmt.Errorf("no email found for vendor ID %s", vendorId)
	}

	return vendor.Email, nil
}
