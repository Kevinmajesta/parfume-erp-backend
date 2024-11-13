package handler

import (
	"net/http"
	"strconv"

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

	NewMo := entity.NewMos("", input.ProductId, input.BomId, input.Qtytoproduce)
	mo, err := h.moService.CreateMo(NewMo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully created a new Manufacture Order", mo))
}

func (h *MoHandler) UpdateMoStatus(c echo.Context) error {
	var input binder.UpdateMoStatusRequest

	// Bind the input
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Input binding error"))
	}

	// Validate input
	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Validation Error: "+err.Error()))
	}

	// Call service to update the manufacture order status
	updatedMo, err := h.moService.UpdateMoStatus(input.MoId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	// Return the updated manufacture order
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully updated manufacture order status", updatedMo))
}

func (h *MoHandler) FindAllMos(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // default page if page parameter is invalid
	}

	materials, err := h.moService.FindAllMos(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data Manufacture Order", materials))
}

func (h *MoHandler) GetMoProfile(c echo.Context) error {
	material_ID := c.Param("id_mo")

	material, err := h.moService.GetMoByID(material_ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to get material"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "successfully displays material data", material))
}

func (h *MoHandler) DeleteMo(c echo.Context) error {
	var input binder.MoDeleteRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	isDeleted, err := h.moService.DeleteMo(input.MoId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses delete Mo", isDeleted))
}
