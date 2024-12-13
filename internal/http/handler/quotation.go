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

type QuoHandler struct {
	quoService service.QuoService
}

func NewQuoHandler(quoService service.QuoService) QuoHandler {
	return QuoHandler{quoService: quoService}
}

func hasDuplicateProduct(products []binder.QUOProductRequest) bool {
	seen := make(map[string]struct{})
	for _, product := range products {
		if _, exists := seen[product.ProductId]; exists {
			return true
		}
		seen[product.ProductId] = struct{}{}
	}
	return false
}

func (h *QuoHandler) CreateRfq(c echo.Context) error {
	var input binder.QUOCreateRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBom(http.StatusBadRequest, "Invalid input"))
	}

	// Periksa apakah vendor ID valid
	exists, err := h.quoService.GetCheckIDProduct(input.CostumerId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, err.Error()))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponseBom(http.StatusNotFound, fmt.Sprintf("Costumer with id %s does not exist", input.CostumerId)))
	}

	// Periksa duplikasi produk
	if hasDuplicateProducts(input.Products) {
		return c.JSON(http.StatusConflict, response.ErrorResponseBom(http.StatusConflict, "Duplicate product IDs found in the request"))
	}

	// Validasi material untuk setiap produk
	for _, product := range input.Products {
		materialCheck, err := h.quoService.GetCheckIDMaterial(product.ProductId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, err.Error()))
		}
		if !materialCheck {
			return c.JSON(http.StatusNotFound, response.ErrorResponseBom(http.StatusNotFound, fmt.Sprintf("Product with id %s does not exist", product.ProductId)))
		}
	}

	// Buat RFQ baru
	newRfq := entity.NewQuo("", input.OrderDate, input.Status, input.CostumerId)

	// Proses produk dan tambahkan VendorId dari input
	var products []entity.QuotationsProduct
	for _, product := range input.Products {
		products = append(products, entity.QuotationsProduct{
			ProductId:   product.ProductId,
			ProductName: product.ProductName,
			Quantity:    product.Quantity,
			UnitPrice:   product.UnitPrice,
			Tax:         product.Tax,
			Subtotal:    product.Subtotal,
			CostumerId:  input.CostumerId,
		})
	}

	// Tambahkan produk ke RFQ
	if len(products) > 0 {
		newRfq.Products = products
	}

	newRfq.CreatedAt = time.Now()
	newRfq.UpdatedAt = time.Now()

	// Simpan RFQ ke database
	bom, err := h.quoService.CreateQuo(newRfq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, "Failed to create Quo"))
	}

	return c.JSON(http.StatusOK, response.BOMResponse{
		Meta: response.Meta{
			Code:    200,
			Message: "Successfully input a new Quo",
		},
		DataBom: *bom,
	})
}

func (h *QuoHandler) UpdateQuo(c echo.Context) error {
	var input binder.QUOUpdateRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBom(http.StatusBadRequest, "Invalid input"))
	}

	// Ambil data RFQ lama untuk validasi
	existingRfq, err := h.quoService.FindQuoById(input.QuotationsId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, err.Error()))
	}
	if existingRfq == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponseBom(http.StatusNotFound, fmt.Sprintf("Quo with id %s not found", input.QuotationsId)))
	}

	// Buat RFQ baru untuk pembaruan
	updatedRfq := entity.UpdateQuo(
		input.QuotationsId,
		input.OrderDate,
		input.Status, // Status baru (kosong jika tidak diubah)
		input.CostumerId,
		existingRfq.Status, // Status lama
	)

	// Tambahkan produk baru dengan VendorId yang sudah diset
	var products []entity.QuotationsProduct
	for _, product := range input.Products {
		// Pastikan VendorId tetap diatur
		products = append(products, entity.QuotationsProduct{
			ProductId:   product.ProductId,
			ProductName: product.ProductName,
			Quantity:    product.Quantity,
			UnitPrice:   product.UnitPrice,
			Tax:         product.Tax,
			Subtotal:    product.Subtotal,
			CostumerId:  product.CostumerId,
		})
	}
	updatedRfq.Products = products

	// Perbarui data RFQ
	result, err := h.quoService.UpdateQuo(updatedRfq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, "Failed to update Quo"))
	}

	return c.JSON(http.StatusOK, response.BOMResponse{
		Meta: response.Meta{
			Code:    200,
			Message: "Successfully updated Quo",
		},
		DataBom: *result,
	})
}

func (h *QuoHandler) DeleteRfq(c echo.Context) error {
	var input binder.QuoDeleteRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	isDeleted, err := h.quoService.DeleteRFQ(input.QuotationsId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses delete Quo", isDeleted))
}

func (h *QuoHandler) FindAllQuo(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // default page if page parameter is invalid
	}

	boms, err := h.quoService.FindAllQuo(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data Quo's", boms))
}

func (h *QuoHandler) FindAllQuoBill(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // default page if page parameter is invalid
	}

	boms, err := h.quoService.FindAllQuoBill(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data Quo's", boms))
}

func (h *QuoHandler) UpdateQuoStatus(c echo.Context) error {
	var input binder.UpdateQuoStatusRequest

	// Bind the input
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Input binding error"))
	}

	// Validate input
	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Validation Error: "+err.Error()))
	}

	// Call service to update the manufacture order status
	updatedMo, err := h.quoService.UpdateQuoStatus(input.QuotationsId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	// Return the updated manufacture order
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully updated Quo status", updatedMo))
}

func (h *QuoHandler) GetQuoOverview(c echo.Context) error {
	rfqId := c.Param("id_quotation")

	if rfqId == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBom(http.StatusBadRequest, "Quo ID cannot be empty"))
	}

	overview, err := h.quoService.CalculateOverview(rfqId)
	if err != nil {
		if err.Error() == "RFQ not found" {
			return c.JSON(http.StatusNotFound, response.ErrorResponseBom(http.StatusNotFound, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.BOMResponse{
		Meta: response.Meta{
			Code:    200,
			Message: "Quo overview retrieved successfully",
		},
		DataBom: overview,
	})
}

func (h *QuoHandler) GetCostumerEmailById(c echo.Context) error {
	// Mendapatkan vendorId dari parameter URL
	vendorId := c.Param("id_costumer")

	// Panggil service untuk mengecek apakah email ada untuk vendorId ini
	email, err := h.quoService.GetEmailByCostumerId(vendorId)
	if err != nil {
		// Jika email tidak ditemukan atau ada error, kirimkan response error
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	// Setelah email ditemukan, dapatkan RFQ ID dari parameter query (atau URL) untuk dikirimkan
	rfqId := c.QueryParam("quotation_id")

	// Panggil service untuk mengirimkan email RFQ
	err = h.quoService.SendQuoEmail(rfqId, email)
	if err != nil {
		// Jika terjadi kesalahan dalam mengirimkan email, kirimkan response error
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Jika berhasil mengirimkan email, kirimkan response sukses
	return c.JSON(http.StatusOK, map[string]string{"message": "Quo email successfully sent"})
}
