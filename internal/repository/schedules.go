package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/cache"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SchedulesRepository interface {
	CreateSchedules(schedules *entity.Schedules) (*entity.Schedules, error)
	FindAllSchedule(page int) ([]entity.Schedules, error)
	UpdateSchedule(admin *entity.Schedules) (*entity.Schedules, error)
	CheckScheduleExists(id uuid.UUID) (bool, error)
	FindScheduleByID(id_schedules uuid.UUID) (*entity.Schedules, error)
	DeleteSchedule(schedule *entity.Schedules) (bool, error)
}

type schedulesRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewSchedulesRepository(db *gorm.DB, cacheable cache.Cacheable) *schedulesRepository {
	return &schedulesRepository{db: db, cacheable: cacheable}
}

func (r *schedulesRepository) CreateSchedules(schedules *entity.Schedules) (*entity.Schedules, error) {
	query := r.db
	if err := query.Create(&schedules).Error; err != nil {
		return schedules, err
	}
	r.cacheable.Delete("FindAllSchedule_page_1")
	return schedules, nil
}

func (r *schedulesRepository) FindScheduleByID(id_schedules uuid.UUID) (*entity.Schedules, error) {
	schedule := new(entity.Schedules)
	if err := r.db.Where("id_schedules = ?", id_schedules).Take(schedule).Error; err != nil {
		return schedule, err
	}
	return schedule, nil
}

func (r *schedulesRepository) FindAllSchedule(page int) ([]entity.Schedules, error) {
	var shcedules []entity.Schedules
	key := fmt.Sprintf("FindAllSchedule_page_%d", page)
	const pageSize = 10

	data, _ := r.cacheable.Get(key)
	if data == "" {
		offset := (page - 1) * pageSize
		if err := r.db.Limit(pageSize).Offset(offset).Find(&shcedules).Error; err != nil {
			return shcedules, err
		}
		marshalledschedules, _ := json.Marshal(shcedules)
		err := r.cacheable.Set(key, marshalledschedules, 5*time.Minute)
		if err != nil {
			return shcedules, err
		}
	} else {
		err := json.Unmarshal([]byte(data), &shcedules)
		if err != nil {
			return shcedules, err
		}
	}
	return shcedules, nil
}

func (r *schedulesRepository) UpdateSchedule(schedule *entity.Schedules) (*entity.Schedules, error) {
	// Use map to store fields to be updated.
	fields := make(map[string]interface{})

	// Update fields only if they are not empty.
	if schedule.Title != "" {
		fields["title"] = schedule.Title
	}
	if schedule.Qty_kolam != "" {
		fields["qty_kolam"] = schedule.Qty_kolam
	}
	if schedule.Date_schedules != "" {
		fields["date_schedules"] = schedule.Date_schedules
	}

	// Update the database in one query.
	if err := r.db.Model(schedule).Where("id_schedules = ?", schedule.SchedulesId).Updates(fields).Error; err != nil {
		return schedule, err
	}
	r.cacheable.Delete("FindAllSchedule_page_1")
	return schedule, nil
}

func (r *schedulesRepository) CheckScheduleExists(id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Schedules{}).Where("id_schedules = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *schedulesRepository) DeleteSchedule(schedule *entity.Schedules) (bool, error) {
	if err := r.db.Delete(&entity.Schedules{}, schedule.SchedulesId).Error; err != nil {
		return false, err
	}
	r.cacheable.Delete("FindAllSchedule_page_1")
	return true, nil
}
