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

func (h *MoHandler) DownloadMOPDF(c echo.Context) error {
	// Retrieve MO ID from the URL parameter (e.g., /mo/:moId)
	moId := c.Param("id_mo")

	if moId == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "MO ID cannot be empty",
		})
	}

	// Fetch the Manufacturing Order (MO) details from the repository using MO ID
	mo, err := h.moService.GetMoByID(moId)
	if err != nil {
		// If the MO is not found, return an error response
		if err.Error() == "MO not found" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Generate the PDF using the GenerateMOPDF method
	pdfData, err := h.moService.GenerateMOPDF(mo)
	if err != nil {
		// If an error occurs while generating the PDF, return an internal server error
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Error generating PDF",
		})
	}

	// Serve the PDF to the client
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=mo_"+mo.MoId+".pdf")
	c.Response().Write(pdfData)

	return nil
}
