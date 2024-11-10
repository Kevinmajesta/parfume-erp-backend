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

type MoRepository interface {
	GetLastMo() (string, error)
	CreateMo(mo *entity.Mos) (*entity.Mos, error)
	FindMoByID(moId string) (*entity.Mos, error)
	UpdateMoStatus(mo *entity.Mos) (*entity.Mos, error)
	FindAllMos(page int) ([]entity.Mos, error)
	DeleteMo(mo *entity.Mos) (bool, error)
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
	r.cacheable.Delete("FindAllMo_page_1")
	return mo, nil
}

func (r *moRepository) FindMoByID(moId string) (*entity.Mos, error) {
	var mo entity.Mos
	err := r.db.First(&mo, "id_mo = ?", moId).Error
	if err != nil {
		return nil, err
	}
	return &mo, nil
}
func (r *moRepository) UpdateMoStatus(mo *entity.Mos) (*entity.Mos, error) {
	if err := r.db.Save(mo).Error; err != nil {
		return nil, err
	}
	r.cacheable.Delete("FindAllMo_page_1")
	return mo, nil
}

func (r *moRepository) FindAllMos(page int) ([]entity.Mos, error) {
	var Mos []entity.Mos
	key := fmt.Sprintf("FindAllMos_page_%d", page)
	const pageSize = 100

	data, _ := r.cacheable.Get(key)
	if data == "" {
		offset := (page - 1) * pageSize
		if err := r.db.Limit(pageSize).Offset(offset).Find(&Mos).Error; err != nil {
			return Mos, err
		}
		marshalledMos, _ := json.Marshal(Mos)
		err := r.cacheable.Set(key, marshalledMos, 5*time.Minute)
		if err != nil {
			return Mos, err
		}
	} else {
		err := json.Unmarshal([]byte(data), &Mos)
		if err != nil {
			return Mos, err
		}
	}
	return Mos, nil
}

func (r *moRepository) DeleteMo(mo *entity.Mos) (bool, error) {
	if err := r.db.Unscoped().Delete(mo).Error; err != nil {
		return false, err
	}
	r.cacheable.Delete("FindAllMo_page_1")
	return true, nil
}
