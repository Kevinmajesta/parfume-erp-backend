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

func hasDuplicateMaterials(materials []binder.MaterialRequest) bool {
    seen := make(map[string]struct{})
    for _, material := range materials {
        if _, exists := seen[material.IdMaterial]; exists {
            return true // Duplikasi ditemukan
        }
        seen[material.IdMaterial] = struct{}{}
    }
    return false // Tidak ada duplikasi
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

	if hasDuplicateMaterials(input.Materials) {
        return c.JSON(http.StatusConflict, response.ErrorResponseBom(http.StatusConflict, "Duplicate material IDs found in the request"))
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

func (h *BOMHandler) DeleteBom(c echo.Context) error {
	var input binder.BomDeleteRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	isDeleted, err := h.bomService.DeleteBom(input.BomId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses delete BoM", isDeleted))
}

func (h *BOMHandler) UpdateBOM(c echo.Context) error {
    // Parse the BOM ID from the URL
    bomId := c.Param("id_bom")

    // Bind the input data
    var input binder.BOMUpdateRequest
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, response.ErrorResponseBom(http.StatusBadRequest, "Invalid input"))
    }

    // Validate the existence of the product
    exists, err := h.bomService.GetCheckIDProduct(input.IdProduct)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, err.Error()))
    }
    if !exists {
        return c.JSON(http.StatusNotFound, response.ErrorResponseBom(http.StatusNotFound, fmt.Sprintf("Product with id %s does not exist", input.IdProduct)))
    }


    // Validasi duplikasi id_material
    if hasDuplicateMaterials(input.Materials) {
        return c.JSON(http.StatusConflict, response.ErrorResponseBom(http.StatusConflict, "Duplicate material IDs found in the request"))
    }

    // Prepare BOM entity using the input and the existing BOM ID
    updatedBomEntity := entity.UpdateBOM(bomId, input.IdProduct, input.ProductName, input.ProductReference, input.Quantity)

    // Prepare materials for update
    var updatedMaterials []entity.BomMaterial
    for _, material := range input.Materials {
        updatedMaterials = append(updatedMaterials, entity.BomMaterial{
            IdBomMaterial: material.IdBomMaterial,
            IdMaterial:    material.IdMaterial,
            MaterialName:  material.MaterialName,
            Quantity:      material.Quantity,
            Unit:          material.Unit,
            BomId:         bomId, // Associate the BOM ID
        })
    }
    updatedBomEntity.Materials = updatedMaterials

    // Call the service to update the BOM and materials
    updatedBom, err := h.bomService.UpdateBOM(updatedBomEntity)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, "Failed to update BoM"))
    }

    return c.JSON(http.StatusOK, response.BOMResponse{
        Meta: response.Meta{
            Code:    200,
            Message: "Successfully updated BoM",
        },
        DataBom: *updatedBom,
    })
}

func (h *BOMHandler) GetBOMByID(c echo.Context) error {
    bomId := c.Param("id_bom")

    bom, err := h.bomService.GetBOMByID(bomId)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, err.Error()))
    }
    if bom == nil {
        return c.JSON(http.StatusNotFound, response.ErrorResponseBom(http.StatusNotFound, fmt.Sprintf("BOM with ID %s not found", bomId)))
    }

    return c.JSON(http.StatusOK, response.BOMResponse{
        Meta: response.Meta{
            Code:    200,
            Message: "Successfully retrieved BoM",
        },
        DataBom: *bom,
    })
}


