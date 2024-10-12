package handler

import (
	"net/http"
	"strconv"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/http/binder"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/service"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) AdminHandler {
	return AdminHandler{adminService: adminService}
}

func (h *AdminHandler) LoginAdmin(c echo.Context) error {
	input := binder.AdminLoginRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	admin, err := h.adminService.LoginAdmin(input.Email, input.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "login success", admin))
}

func (h *AdminHandler) FindAllUser(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // default page if page parameter is invalid
	}

	users, err := h.adminService.FindAllUser(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data users", users))
}

func (h *AdminHandler) CreateAdmin(c echo.Context) error {
	input := binder.AdminCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	if h.adminService.EmailExists(input.Email) {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "email is already in use"))
	}

	newAdmin := entity.NewAdmin(input.Fullname, input.Email, input.Password, input.Role, input.Phone, input.Verification)
	admin, err := h.adminService.CreateAdmin(newAdmin)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully created a new admin, the email has been sent", admin))
}

func (h *AdminHandler) UpdateAdmin(c echo.Context) error {
	var input binder.AdminUpdateRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	id, err := uuid.Parse(input.Admin_ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid admin ID"))
	}
	exists, err := h.adminService.CheckUserExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "could not verify user existence"))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "user ID does not exist"))
	}
	inputAdmin := entity.UpdateAdmin(id, input.Fullname, input.Email, input.Password, input.Role, input.Phone, input.Verification)

	updatedAdmin, err := h.adminService.UpdateAdmin(inputAdmin)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success update admin", updatedAdmin))
}

func (h *AdminHandler) DeleteAdmin(c echo.Context) error {
	var input binder.AdminDeleteRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	id := uuid.MustParse(input.Admin_ID)

	isDeleted, err := h.adminService.DeleteAdmin(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success delete admin", isDeleted))
}
