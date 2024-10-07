package service

import (
	// Import log package

	"errors"

	"github.com/Kevinmajesta/webPemancingan/internal/entity"
	"github.com/Kevinmajesta/webPemancingan/internal/repository"
	"github.com/google/uuid"
)

type SchedulesService interface {
	CreateSchedules(admin *entity.Schedules) (*entity.Schedules, error)
	FindAllSchedule(page int) ([]entity.Schedules, error)
	CheckScheduleExists(id uuid.UUID) (bool, error)
	UpdateSchedule(schedule *entity.Schedules) (*entity.Schedules, error)
	DeleteSchedule(id_schedules uuid.UUID) (bool, error)
}

type schedulesService struct {
	schedulesRepository repository.SchedulesRepository
}

func NewSchedulesService(schedulesRepository repository.SchedulesRepository) *schedulesService {
	return &schedulesService{
		schedulesRepository: schedulesRepository,
	}
}

func (s *schedulesService) CreateSchedules(schedules *entity.Schedules) (*entity.Schedules, error) {
	return s.schedulesRepository.CreateSchedules(schedules)
}

func (s *schedulesService) FindAllSchedule(page int) ([]entity.Schedules, error) {
	schedule, err := s.schedulesRepository.FindAllSchedule(page)
	if err != nil {
		return nil, err
	}

	formattedSchedule := make([]entity.Schedules, 0)
	for _, v := range schedule {
		formattedSchedule = append(formattedSchedule, v)
	}

	return formattedSchedule, nil
}

func (s *schedulesService) UpdateSchedule(schedule *entity.Schedules) (*entity.Schedules, error) {
	if schedule.Title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if schedule.Qty_kolam == "" {
		return nil, errors.New("quantity cannot be empty")
	}
	if schedule.Date_schedules == "" {
		return nil, errors.New("date cannot be empty")
	}

	return s.schedulesRepository.UpdateSchedule(schedule)
}

func (s *schedulesService) DeleteSchedule(id_schedules uuid.UUID) (bool, error) {
	schedule, err := s.schedulesRepository.FindScheduleByID(id_schedules)
	if err != nil {
		return false, err
	}

	return s.schedulesRepository.DeleteSchedule(schedule)
}

func (s *schedulesService) CheckScheduleExists(id uuid.UUID) (bool, error) {
	return s.schedulesRepository.CheckScheduleExists(id)
}
