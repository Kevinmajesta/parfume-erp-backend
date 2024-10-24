package handler

import (
	"net/http"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/http/binder"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/service"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/response"
	"github.com/labstack/echo/v4"
)

type MoHandler struct {
	moService service.MoService
}

func NewMoHandler(moService service.MoService) MoHandler {
	return MoHandler{moService: moService}
}

func (h *MoHandler) CreateMo(c echo.Context) error {
	input := binder.MoCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	NewMo := entity.NewMos("",input.BomId, input.ProductId, input.Qtytoproduce)
	mo, err := h.moService.CreateMo(NewMo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully created a new Manufacture Order", mo))
}
