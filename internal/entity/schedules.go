package entity

import (
	"github.com/google/uuid"
)

type Schedules struct {
	SchedulesId    uuid.UUID `gorm:"column:id_schedules;primary_key"`
	Title          string    `json:"title"`
	Qty_kolam      string    `json:"qty_kolam"`
	Date_schedules string    `json:"date_schedules"`
	Auditable
}

func NewSchedules(title, qty_kolam, date_schedules string) *Schedules {
	return &Schedules{
		SchedulesId:    uuid.New(),
		Title:          title,
		Qty_kolam:      qty_kolam,
		Date_schedules: date_schedules,
		Auditable:      NewAuditable(),
	}
}

func UpdateSchedule(id_schedules uuid.UUID, title, qty_kolam, date_schedules string) *Schedules {
	return &Schedules{
		SchedulesId:    id_schedules,
		Title:          title,
		Qty_kolam:      qty_kolam,
		Date_schedules: date_schedules,
		Auditable:      UpdateAuditable(),
	}
}
