package repository

import (
	"errors"

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
	return rfq, nil
}

func (r *rfqRepository) GetRfqById(rfqId string) (*entity.Rfqs, error) {
	var rfq entity.Rfqs
	if err := r.db.Where("id_rfq = ?", rfqId).First(&rfq).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rfq, nil
}


