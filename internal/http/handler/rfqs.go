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

func (h *RfqHandler) UpdateRfqStatus(c echo.Context) error {
	var input binder.UpdateRfqStatusRequest

	// Bind the input
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Input binding error"))
	}

	// Validate input
	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Validation Error: "+err.Error()))
	}

	// Call service to update the manufacture order status
	updatedMo, err := h.rfqService.UpdateRfqStatus(input.RfqId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	// Return the updated manufacture order
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully updated RFQ status", updatedMo))
}

func (h *RfqHandler) FindAllRfq(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // default page if page parameter is invalid
	}

	boms, err := h.rfqService.FindAllRfq(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data RFQ's", boms))
}

func (h *RfqHandler) FindAllRfqBill(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // default page if page parameter is invalid
	}

	boms, err := h.rfqService.FindAllRfqBill(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data RFQ's", boms))
}

func (h *RfqHandler) GetRfqOverview(c echo.Context) error {
	rfqId := c.Param("id_rfq")

	if rfqId == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBom(http.StatusBadRequest, "RFQ ID cannot be empty"))
	}

	overview, err := h.rfqService.CalculateOverview(rfqId)
	if err != nil {
		if err.Error() == "RFQ not found" {
			return c.JSON(http.StatusNotFound, response.ErrorResponseBom(http.StatusNotFound, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ErrorResponseBom(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.BOMResponse{
		Meta: response.Meta{
			Code:    200,
			Message: "RFQ overview retrieved successfully",
		},
		DataBom: overview,
	})
}

func (h *RfqHandler) GetVendorEmailById(c echo.Context) error {
	// Mendapatkan vendorId dari parameter URL
	vendorId := c.Param("id_vendor")

	// Panggil service untuk mengecek apakah email ada untuk vendorId ini
	email, err := h.rfqService.GetEmailByVendorId(vendorId)
	if err != nil {
		// Jika email tidak ditemukan atau ada error, kirimkan response error
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	// Setelah email ditemukan, dapatkan RFQ ID dari parameter query (atau URL) untuk dikirimkan
	rfqId := c.QueryParam("rfq_id")

	// Panggil service untuk mengirimkan email RFQ
	err = h.rfqService.SendRfqEmail(rfqId, email)
	if err != nil {
		// Jika terjadi kesalahan dalam mengirimkan email, kirimkan response error
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Jika berhasil mengirimkan email, kirimkan response sukses
	return c.JSON(http.StatusOK, map[string]string{"message": "RFQ email successfully sent"})
}

func (h *RfqHandler) DeleteRfq(c echo.Context) error {
	var input binder.RFQDeleteRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	isDeleted, err := h.rfqService.DeleteRFQ(input.RfqId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses delete RFQ", isDeleted))
}

func (h *RfqHandler) HandleCreateRfqPDF(c echo.Context) error {
	// Get the RFQ ID from the URL parameters
	rfqId := c.Param("id_rfq")

	// Call the service method to create the PDF
	pdfBytes, err := h.rfqService.CreateRfqPDF(rfqId, "")
	if err != nil {
		// If an error occurs while generating the PDF, return error response
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Set the appropriate headers for PDF download
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=rfq-"+rfqId+".pdf")
	c.Response().Header().Set("Content-Length", string(len(pdfBytes)))

	// Write the PDF content to the response
	if _, err := c.Response().Write(pdfBytes); err != nil {
		// Return error if there is an issue writing the PDF content
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return nil
}
