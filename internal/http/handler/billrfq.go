package handler

import (
	"net/http"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/http/binder"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/service"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/response"
	"github.com/labstack/echo/v4"
)

type BillrfqHandler struct {
	billrfqService service.BillrfqService
}

func NewBillrfqHandler(billrfqService service.BillrfqService) BillrfqHandler {
	return BillrfqHandler{billrfqService: billrfqService}
}

func (h *BillrfqHandler) CreateMo(c echo.Context) error {
	input := binder.BillRfqCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	NewMo := entity.NewBillrfq("", input.VendorId, input.Bill_date, input.Payment)
	mo, err := h.billrfqService.CreateBill(NewMo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully created a new bill", mo))
}
