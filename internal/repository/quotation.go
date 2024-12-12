package repository

import (
	"errors"

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
