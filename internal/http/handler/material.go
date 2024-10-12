package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/http/binder"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/service"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MaterialHandler struct {
	materialService service.MaterialService
}

func NewMaterialHandler(materialService service.MaterialService) MaterialHandler {
	return MaterialHandler{materialService: materialService}
}

func (h *MaterialHandler) CreateMaterial(c echo.Context) error {
	input := binder.MaterialCreateRequest{}

	// Bind input JSON ke struct
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "There is an input error"))
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Failed to retrieve image"))
	}

	// Check image format
	chckFormat := strings.ToLower(filepath.Ext(file.Filename))
	if chckFormat != ".jpg" && chckFormat != ".jpeg" && chckFormat != ".png" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid image format. Only jpg, jpeg, and png are allowed"))
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to open image"))
	}
	defer src.Close()

	imageID := uuid.New()
	imageFilename := fmt.Sprintf("%s%s", imageID, filepath.Ext(file.Filename))
	imagePath := filepath.Join("assets", "images", imageFilename)

	dst, err := os.Create(imagePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to create image file"))
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to copy image file"))
	}

	// Buat entitas produk dari input
	newMaterial := &entity.Materials{
		Materialname:     input.MaterialName,
		Materialcategory: input.MaterialCategory,
		Sellprice:        input.SellPrice,
		Makeprice:        input.MakePrice,
		Unit:             input.Unit,
		Description:      input.Description,
		Image:            "/assets/images/" + imageFilename,
	}

	// Panggil service untuk membuat produk baru
	material, err := h.materialService.CreateMaterial(newMaterial)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Response jika berhasil
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully input a new material", material))
}

func (h *MaterialHandler) UpdateMaterial(c echo.Context) error {
	var input binder.MaterialUpdateRequest
	var imageURL string

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	if input.MaterialId == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Material ID cannot be empty"))
	}
	exists, err := h.materialService.CheckMaterialExists(input.MaterialId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "could not verify produk existence"))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "material ID does not exist"))
	}

	file, err := c.FormFile("image")

	if err == nil {
		// Check image format
		chckFormat := strings.ToLower(filepath.Ext(file.Filename))
		if chckFormat != ".jpg" && chckFormat != ".jpeg" && chckFormat != ".png" {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid image format. Only jpg, jpeg, and png are allowed"))
		}

		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to open image"))
		}
		defer src.Close()

		imageID := uuid.New()
		imageFilename := fmt.Sprintf("%s%s", imageID, filepath.Ext(file.Filename))
		imagePath := filepath.Join("assets", "images", imageFilename)

		dst, err := os.Create(imagePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to create image file"))
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to copy image file"))
		}

		imageURL = "/assets/images/" + imageFilename
	} else {
		imageURL = ""
	}

	inputMaterial := entity.UpdateMaterials(input.MaterialId, input.MaterialtName, input.MaterialtCategory, input.SellPrice, input.MakePrice, input.Unit, input.Description, imageURL)

	updatedMaterial, err := h.materialService.UpdateMaterial(inputMaterial)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success update material", updatedMaterial))
}

func (h *MaterialHandler) DeleteMaterial(c echo.Context) error {
	var input binder.MaterialDeleteRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	isDeleted, err := h.materialService.DeleteMaterial(input.MaterialId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses delete material", isDeleted))
}

func (h *MaterialHandler) FindAllMaterial(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // default page if page parameter is invalid
	}

	materials, err := h.materialService.FindAllMaterial(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data materials", materials))
}

func (h *MaterialHandler) SearchMaterials(c echo.Context) error {
	name := c.QueryParam("materialname")
	materials, err := h.materialService.SearchMaterialsByName(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// Check if insert is empty
	if name == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Name required"))
	}
	// Check if title is not available
	if len(materials) == 0 {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "Material not found"))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data Material", materials))
}