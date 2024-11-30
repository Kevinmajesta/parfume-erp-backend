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

type RfqHandler struct {
	rfqService service.RfqService
}

func NewRfqHandler(rfqService service.RfqService) RfqHandler {
	return RfqHandler{rfqService: rfqService}
}

func hasDuplicateProducts(products []binder.ProductRequest) bool {
	seen := make(map[string]struct{})
	for _, product := range products {
		if _, exists := seen[product.ProductId]; exists {
			return true // Duplikasi ditemukan
		}
		seen[product.ProductId] = struct{}{}
	}
	return false // Tidak ada duplikasi
}

func (h *RfqHandler) CreateRfq(c echo.Context) error {
	var input binder.RFQCreateRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBom(http.StatusBadRequest, "Invalid input"))
	}

	// Periksa apakah vendor ID valid
	exists, err := h.rfqService.GetCheckIDProduct(input.VendorId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, err.Error()))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponseBom(http.StatusNotFound, fmt.Sprintf("Vendor with id %s does not exist", input.VendorId)))
	}

	// Periksa duplikasi produk
	if hasDuplicateProducts(input.Products) {
		return c.JSON(http.StatusConflict, response.ErrorResponseBom(http.StatusConflict, "Duplicate product IDs found in the request"))
	}

	// Validasi material untuk setiap produk
	for _, product := range input.Products {
		materialCheck, err := h.rfqService.GetCheckIDMaterial(product.ProductId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, err.Error()))
		}
		if !materialCheck {
			return c.JSON(http.StatusNotFound, response.ErrorResponseBom(http.StatusNotFound, fmt.Sprintf("Product with id %s does not exist", product.ProductId)))
		}
	}

	// Buat RFQ baru
	newRfq := entity.NewRfqs("", input.OrderDate, input.Status, input.VendorId)

	// Proses produk dan tambahkan VendorId dari input
	var products []entity.RfqsProduct
	for _, product := range input.Products {
		products = append(products, entity.RfqsProduct{
			ProductId:   product.ProductId,
			ProductName: product.ProductName,
			Quantity:    product.Quantity,
			UnitPrice:   product.UnitPrice,
			Tax:         product.Tax,
			Subtotal:    product.Subtotal,
			VendorId:    input.VendorId, // Tambahkan VendorId
		})
	}

	// Tambahkan produk ke RFQ
	if len(products) > 0 {
		newRfq.Products = products
	}

	newRfq.CreatedAt = time.Now()
	newRfq.UpdatedAt = time.Now()

	// Simpan RFQ ke database
	bom, err := h.rfqService.CreateRfq(newRfq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, "Failed to create RFQ"))
	}

	return c.JSON(http.StatusOK, response.BOMResponse{
		Meta: response.Meta{
			Code:    200,
			Message: "Successfully input a new RFQ",
		},
		DataBom: *bom,
	})
}

func (h *RfqHandler) UpdateRfq(c echo.Context) error {
	var input binder.RFQUpdateRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBom(http.StatusBadRequest, "Invalid input"))
	}

	// Ambil data RFQ lama untuk validasi
	existingRfq, err := h.rfqService.FindRfqById(input.RfqId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, err.Error()))
	}
	if existingRfq == nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponseBom(http.StatusNotFound, fmt.Sprintf("RFQ with id %s not found", input.RfqId)))
	}

	// Buat RFQ baru untuk pembaruan
	updatedRfq := entity.UpdateRfqs(
		input.RfqId,
		input.OrderDate,
		input.Status, // Status baru (kosong jika tidak diubah)
		input.VendorId,
		existingRfq.Status, // Status lama
	)

	// Tambahkan produk baru dengan VendorId yang sudah diset
	var products []entity.RfqsProduct
	for _, product := range input.Products {
		// Pastikan VendorId tetap diatur
		products = append(products, entity.RfqsProduct{
			ProductId:   product.ProductId,
			ProductName: product.ProductName,
			Quantity:    product.Quantity,
			UnitPrice:   product.UnitPrice,
			Tax:         product.Tax,
			Subtotal:    product.Subtotal,
			VendorId:    product.VendorId,
		})
	}
	updatedRfq.Products = products

	// Perbarui data RFQ
	result, err := h.rfqService.UpdateRfq(updatedRfq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, "Failed to update RFQ"))
	}

	return c.JSON(http.StatusOK, response.BOMResponse{
		Meta: response.Meta{
			Code:    200,
			Message: "Successfully updated RFQ",
		},
		DataBom: *result,
	})
}
