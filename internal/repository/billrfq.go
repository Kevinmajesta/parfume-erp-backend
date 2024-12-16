package repository

import (
	"errors"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/cache"
	"gorm.io/gorm"
)

type BillrfqRepository interface {
	GetLastMo() (string, error)
	CreateMo(mo *entity.Billrfq) (*entity.Billrfq, error)
}

type billrfqRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewBillrfqRepository(db *gorm.DB, cacheable cache.Cacheable) *billrfqRepository {
	return &billrfqRepository{db: db, cacheable: cacheable}
}

func (r *billrfqRepository) GetLastMo() (string, error) {
	var lastMo entity.Billrfq
	err := r.db.Order("id_bill DESC").First(&lastMo).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "BRQ-00000", nil
	} else if err != nil {
		return "", err
	}

	return lastMo.BillrfqId, nil
}

func (r *billrfqRepository) CreateMo(mo *entity.Billrfq) (*entity.Billrfq, error) {
	if err := r.db.Create(&mo).Error; err != nil {
		return mo, err
	}
	r.cacheable.Delete("FindAllMo_page_1")
	return mo, nil
}
