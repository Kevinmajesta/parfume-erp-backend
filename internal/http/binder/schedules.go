package binder

type SchedulesCreateRequest struct {
	Title          string `json:"title" validate:"required"`
	Qty_kolam      string `json:"qty_kolam" validate:"required,email"`
	Date_schedules string `json:"date_schedules" validate:"required"`
}

type SchedulesUpdateRequest struct {
	Schedule_ID    string `param:"id_schedules" validate:"required"`
	Title          string `json:"title" validate:"required"`
	Qty_kolam      string `json:"qty_kolam" validate:"required,email"`
	Date_schedules string `json:"date_schedules" validate:"required"`
}

type SchedulesDeleteRequest struct {
	Schedule_ID string `param:"id_schedules" validate:"required"`
}
