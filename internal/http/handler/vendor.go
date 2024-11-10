package handler

import (
	"net/http"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/http/binder"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/service"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/response"
	"github.com/labstack/echo/v4"
)

type VendorHandler struct {
	vendorService service.VendorService
}

func NewVendorHandler(vendorService service.VendorService) VendorHandler {
	return VendorHandler{vendorService: vendorService}
}

func (h *VendorHandler) CreateVendor(c echo.Context) error {
	input := binder.VendorCreateRequest{}

	// Bind input JSON ke struct
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "There is an input error"))
	}

	// Buat entitas produk dari input
	newVendor := &entity.Vendors{
		Vendorname: input.Vendorname,
		Addressone: input.Addressone,
		Addresstwo: input.Addresstwo,
		Phone:      input.Phone,
		Email:      input.Email,
		Website:    input.Website,
	}

	// Panggil service untuk membuat produk baru
	vendor, err := h.vendorService.CreateVendor(newVendor)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Response jika berhasil
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully input a new vendor", vendor))
}

func (h *VendorHandler) UpdateVendor(c echo.Context) error {
	var input binder.VendorUpdateRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	if input.VendorId == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Vendor ID cannot be empty"))
	}
	exists, err := h.vendorService.CheckVendorExists(input.VendorId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "could not verify vendor existence"))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "vendor ID does not exist"))
	}

	inputUser := entity.UpdateVendor(input.VendorId, input.Vendorname, input.Addressone, input.Addresstwo, input.Phone, input.Email, input.Website)

	updatedProduk, err := h.vendorService.UpdateVendor(inputUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success update vendor", updatedProduk))
}

func (h *VendorHandler) DeleteVendor(c echo.Context) error {
	var input binder.VendorDeleteRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	isDeleted, err := h.vendorService.DeleteVendor(input.VendorId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses delete vendor", isDeleted))
}
