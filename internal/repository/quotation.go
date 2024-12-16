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

type QuoRepository interface {
	GetLastQuo() (string, error)
	CreateQuo(rfq *entity.Quotations) (*entity.Quotations, error)
	CheckProductExists(productId string) (bool, error)
	GetQuoById(rfqId string) (*entity.Quotations, error)
	UpdateQuo(rfq *entity.Quotations) (*entity.Quotations, error)
	DeleteQuo(mo *entity.Quotations) (bool, error)
	FindAllQuo(page int) ([]entity.Quotations, error)
	FindAllQuoBill(page int) ([]entity.Quotations, error)
	UpdateQuoStatus(rfq *entity.Quotations) (*entity.Quotations, error)
	CheckEmailExistsByCostumerId(vendorId string) (string, error)
}

type quoRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewQuoRepository(db *gorm.DB, cacheable cache.Cacheable) *quoRepository {
	return &quoRepository{db: db, cacheable: cacheable}
}

func (r *quoRepository) GetLastQuo() (string, error) {
	var lastRfq entity.Quotations
	err := r.db.Order("id_quotation DESC").First(&lastRfq).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "R-00000", nil
	} else if err != nil {
		return "", err
	}

	return lastRfq.QuotationsId, nil
}

func (r *quoRepository) CreateQuo(rfq *entity.Quotations) (*entity.Quotations, error) {
	if err := r.db.Create(&rfq).Error; err != nil {
		return rfq, err
	}
	r.cacheable.Delete("FindAllQuo_page_1")
	r.cacheable.Delete("FindAllQuoBill_page_1")
	return rfq, nil
}

func (r *quoRepository) CheckProductExists(productId string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Costumers{}).Where("id_costumer = ?", productId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *quoRepository) GetQuoById(rfqId string) (*entity.Quotations, error) {
	var rfq entity.Quotations
	if err := r.db.Where("id_quotation = ?", rfqId).
		Preload("Products").
		First(&rfq).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rfq, nil
}

func (r *quoRepository) UpdateQuo(rfq *entity.Quotations) (*entity.Quotations, error) {
	if err := r.db.Model(&entity.Quotations{}).Where("id_quotation = ?", rfq.QuotationsId).Updates(rfq).Error; err != nil {

		return nil, err
	}
	r.cacheable.Delete("FindAllQuo_page_1")
	r.cacheable.Delete("FindAllQuoBill_page_1")
	return rfq, nil
}

func (r *quoRepository) DeleteQuo(mo *entity.Quotations) (bool, error) {
	if err := r.db.Unscoped().Delete(mo).Error; err != nil {
		return false, err
	}
	r.cacheable.Delete("FindAllQuoBill_page_1")
	r.cacheable.Delete("FindAllQuo_page_1")
	return true, nil
}

func (r *quoRepository) FindAllQuo(page int) ([]entity.Quotations, error) {
	var Rfq []entity.Quotations
	key := fmt.Sprintf("FindAllQuo_page_%d", page)
	const pageSize = 100

	data, _ := r.cacheable.Get(key)
	if data == "" {
		offset := (page - 1) * pageSize
		if err := r.db.Where("status IN ?", []string{"QUOTATION", "Sales Order"}).
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
	r.cacheable.Delete("FindAllQuo_page_1")
	r.cacheable.Delete("FindAllQuoBill_page_1")
	return Rfq, nil
}

func (r *quoRepository) FindAllQuoBill(page int) ([]entity.Quotations, error) {
	var Rfq []entity.Quotations
	key := fmt.Sprintf("FindAllQuoBill_page_%d", page)
	const pageSize = 100

	data, _ := r.cacheable.Get(key)
	if data == "" {
		offset := (page - 1) * pageSize
		if err := r.db.Where("status IN ?", []string{"Sales Order", "Invoiced", "Delivery", "Done"}).
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
	r.cacheable.Delete("FindAllQuo_page_1")
	r.cacheable.Delete("FindAllQuoBill_page_1")
	return Rfq, nil
}

func (r *quoRepository) UpdateQuoStatus(rfq *entity.Quotations) (*entity.Quotations, error) {
	if err := r.db.Model(&entity.Quotations{}).
		Where("id_quotation = ?", rfq.QuotationsId).
		Updates(map[string]interface{}{
			"order_date":  rfq.OrderDate,
			"id_costumer": rfq.CostumerId,
			"status":      rfq.Status,
			"updated_at":  rfq.UpdatedAt,
		}).Error; err != nil {
		return nil, err
	}

	r.cacheable.Delete("FindAllQuo_page_1")
	r.cacheable.Delete("FindAllQuoBill_page_1")
	return rfq, nil
}

func (r *quoRepository) CheckEmailExistsByCostumerId(vendorId string) (string, error) {
	var vendor entity.Costumers
	err := r.db.Table("costumers").Where("id_costumer = ?", vendorId).Select("email").Scan(&vendor).Error
	if err != nil {
		return "", err
	}

	if vendor.Email == "" {
		return "", fmt.Errorf("no email found for costuemr ID %s", vendorId)
	}

	return vendor.Email, nil
}
