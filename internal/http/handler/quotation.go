package handler

import (
	"fmt"
	"net/http"
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
