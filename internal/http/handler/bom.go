package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/http/binder"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/service"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/response"
	"github.com/labstack/echo/v4"
)

type BOMHandler struct {
	bomService service.BOMService
}

func NewBOMHandler(bomService service.BOMService) *BOMHandler {
	return &BOMHandler{bomService: bomService}
}

func generateMaterialId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "MTR-%d", &newNumber)
		newNumber++
	}
	return fmt.Sprintf("MTR-%05d", newNumber)
}

func (h *BOMHandler) CreateBOM(c echo.Context) error {
	var input binder.BOMCreateRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBom(http.StatusBadRequest, "Invalid input"))
	}

	exists, err := h.bomService.GetCheckIDProduct(input.IdProduct)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, err.Error()))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponseBom(http.StatusNotFound, fmt.Sprintf("Product with id %s does not exist", input.IdProduct)))
	}

	for _, material := range input.Materials {
		fmt.Println("Received Material Id:", material.IdMaterial) // Log ID material

		materialcheck, err := h.bomService.GetCheckIDMaterial(material.IdMaterial)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, err.Error()))
		}
		if !materialcheck {
			return c.JSON(http.StatusNotFound, response.ErrorResponseBom(http.StatusNotFound, fmt.Sprintf("Material with id %s does not exist", material.IdMaterial)))
		}
	}

	newBom := entity.NewBom("", input.IdProduct, input.ProductName, input.ProductReference, input.Quantity)

	var materials []entity.BomMaterial
	for _, material := range input.Materials {
		materials = append(materials, entity.BomMaterial{
			IdMaterial:   material.IdMaterial,
			MaterialName: material.MaterialName,
			Quantity:     material.Quantity,
			Unit:         material.Unit,
		})
	}

	if len(materials) > 0 {
		newBom.Materials = materials
	}

	newBom.CreatedAt = time.Now()
	newBom.UpdatedAt = time.Now()

	bom, err := h.bomService.CreateBOM(newBom)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, "Failed to create BoM"))
	}

	return c.JSON(http.StatusOK, response.BOMResponse{
		Meta: response.Meta{
			Code:    200,
			Message: "Successfully input a new BoM",
		},
		DataBom: *bom,
	})
}

func (h *BOMHandler) FindAllBom(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // default page if page parameter is invalid
	}

	boms, err := h.bomService.FindAllBom(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data materials", boms))
}
