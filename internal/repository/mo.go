package repository

import (
	"errors"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/cache"
	"gorm.io/gorm"
)

type MoRepository interface {
	GetLastMo() (string, error)
	CreateMo(mo *entity.Mos) (*entity.Mos, error)
}

type moRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewMoRepository(db *gorm.DB, cacheable cache.Cacheable) *moRepository {
	return &moRepository{db: db, cacheable: cacheable}
}

func (r *moRepository) GetLastMo() (string, error) {
	var lastMo entity.Mos
	err := r.db.Order("id_mo DESC").First(&lastMo).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "MO-00000", nil
	} else if err != nil {
		return "", err
	}

	return lastMo.MoId, nil
}

func (r *moRepository) CreateMo(mo *entity.Mos) (*entity.Mos, error) {
	if err := r.db.Create(&mo).Error; err != nil {
		return mo, err
	}
	// r.cacheable.Delete("FindAllMos_page_1")
	// r.cacheable.Delete("FindAllMos_page_2")
	return mo, nil
}