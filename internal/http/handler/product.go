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

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) ProductHandler {
	return ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	input := binder.ProductCreateRequest{}

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
	newProduct := &entity.Products{
		Productname:     input.ProductName,
		Productcategory: input.ProductCategory,
		Sellprice:       input.SellPrice,
		Makeprice:       input.MakePrice,
		Pajak:           input.Pajak,
		Description:     input.Description,
		Image:           "/assets/images/" + imageFilename,
	}

	// Panggil service untuk membuat produk baru
	product, err := h.productService.CreateProduct(newProduct)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Response jika berhasil
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully input a new product", product))
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	var input binder.ProductUpdateRequest
	var imageURL string

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	if input.ProdukId == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Product ID cannot be empty"))
	}
	exists, err := h.productService.CheckProductExists(input.ProdukId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "could not verify produk existence"))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "produk ID does not exist"))
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

	inputUser := entity.UpdateProduct(input.ProdukId, input.ProductName, input.ProductCategory, input.SellPrice, input.MakePrice, input.Pajak, input.Description, imageURL)

	updatedProduk, err := h.productService.UpdateProduct(inputUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success update produk", updatedProduk))
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	var input binder.ProdukDeleteRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	isDeleted, err := h.productService.DeleteProduct(input.ProdukId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses delete product", isDeleted))
}

func (h *ProductHandler) FindAllProduct(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // default page if page parameter is invalid
	}

	products, err := h.productService.FindAllProduct(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data products", products))
}

func (h *ProductHandler) SearchProducts(c echo.Context) error {
	name := c.QueryParam("productname")
	products, err := h.productService.SearchProductsByName(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// Check if insert is empty
	if name == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Name required"))
	}
	// Check if title is not available
	if len(products) == 0 {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "Product not found"))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data product", products))
}

func (h *ProductHandler) DownloadProductPDF(c echo.Context) error {
	id := c.Param("id_product")

	// Panggil service untuk generate PDF
	fileName, err := h.productService.GenerateProductPDF(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	// Kirim file PDF ke client untuk di-download
	return c.File(fileName)
}

// GenerateBarcodePDFHandler untuk mengunduh PDF barcode
func (h *ProductHandler) GenerateBarcodePDFHandler(c echo.Context) error {
	id := c.Param("id_product")

	fileName, err := h.productService.GenerateBarcode(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.File(fileName)
}

func (h *ProductHandler) GenerateAllProductsPDFHandler(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // default page if page parameter is invalid
	}

	fileName, err := h.productService.GenerateAllProductsPDF(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.File(fileName)
}
