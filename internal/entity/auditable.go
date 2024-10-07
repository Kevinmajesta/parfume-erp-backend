package entity

import (
	"time"

	"gorm.io/gorm"
)

type Auditable struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func NewAuditable() Auditable {
	return Auditable{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}
}

func UpdateAuditable() Auditable {
	return Auditable{
		UpdatedAt: time.Now(),
	}
}
