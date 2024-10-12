package handler

import (
	"net/http"
	"strconv"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/http/binder"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/service"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SchedulesHandler struct {
	schedulesService service.SchedulesService
}

func NewSchedulesHandler(schedulesService service.SchedulesService) SchedulesHandler {
	return SchedulesHandler{schedulesService: schedulesService}
}

func (h *SchedulesHandler) CreateSchedules(c echo.Context) error {
	input := binder.SchedulesCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	NewSchedules := entity.NewSchedules(input.Title, input.Qty_kolam, input.Date_schedules)
	schedule, err := h.schedulesService.CreateSchedules(NewSchedules)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully created a new schedules", schedule))
}

func (h *SchedulesHandler) FindAllSchedule(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // default page if page parameter is invalid
	}

	users, err := h.schedulesService.FindAllSchedule(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data schedules", users))
}

func (h *SchedulesHandler) UpdateSchedule(c echo.Context) error {
	var input binder.SchedulesUpdateRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	id, err := uuid.Parse(input.Schedule_ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid schedule ID"))
	}
	exists, err := h.schedulesService.CheckScheduleExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "could not verify schedule existence"))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "schedule ID does not exist"))
	}
	inputSchedule := entity.UpdateSchedule(id, input.Title, input.Qty_kolam, input.Date_schedules)

	updatedSchedule, err := h.schedulesService.UpdateSchedule(inputSchedule)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success update admin", updatedSchedule))
}

func (h *SchedulesHandler) DeleteSchedule(c echo.Context) error {
	var input binder.SchedulesDeleteRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	id := uuid.MustParse(input.Schedule_ID)

	isDeleted, err := h.schedulesService.DeleteSchedule(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success delete schedule", isDeleted))
}
