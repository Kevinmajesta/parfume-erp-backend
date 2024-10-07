package handler

import (
	"net/http"

	"github.com/Kevinmajesta/webPemancingan/internal/entity"
	"github.com/Kevinmajesta/webPemancingan/internal/http/binder"
	"github.com/Kevinmajesta/webPemancingan/internal/service"
	"github.com/Kevinmajesta/webPemancingan/pkg/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{userService: userService}
}

func (h *UserHandler) LoginUser(c echo.Context) error {
	input := new(binder.UserLoginRequest)

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	user, err := h.userService.LoginUser(input.Email, input.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "login success", user))
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	input := binder.UserCreateRequest{}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	if h.userService.EmailExists(input.Email) {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "email is already in use"))
	}

	newUser := entity.NewUser(input.Fullname, input.Email, input.Password, input.Phone, input.Role, input.Status, input.Verification)
	user, err := h.userService.CreateUser(newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully created a new user, the email has been sent", user))
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	var input binder.UserUpdateRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	id, err := uuid.Parse(input.User_ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "invalid user ID"))
	}

	exists, err := h.userService.CheckUserExists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "could not verify user existence"))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "user ID does not exist"))
	}

	inputUser := entity.UpdateUser(id, input.Fullname, input.Email, input.Password, input.Phone, input.Role, input.Status, input.Verification)

	updatedUser, err := h.userService.UpdateUser(inputUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success update user", updatedUser))
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	var input binder.UserDeleteRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	id := uuid.MustParse(input.User_ID)

	isDeleted, err := h.userService.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success delete user", isDeleted))
}

func (h *UserHandler) GetUserProfile(c echo.Context) error {
	// Dapatkan ID pengguna dari parameter URL
	user_ID := c.Param("id_user")

	// Panggil layanan untuk mendapatkan profil pengguna berdasarkan ID
	user, err := h.userService.GetUserProfileByID(user_ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to get user profile"))
	}

	// Mengembalikan data profil pengguna sebagai respons JSON
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "successfully displays user data", user))
}

// func (h *UserHandler) RequestPasswordReset(c echo.Context) error {
// 	var req binder.PasswordResetRequest
// 	err := json.NewDecoder(c.Request().Body).Decode(&req)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid request"))
// 	}

// 	err = h.userService.RequestPasswordReset(req.Email)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusBadRequest, "Invalid request"))
// 	}

// 	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Reset code sent", nil))
// }

func (h *UserHandler) ResetPassword(c echo.Context) error {
	var req binder.ResetPassword
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	if err := h.userService.ResetPassword(req.ResetCode, req.Password); err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success reset password ", nil))
}

func (h *UserHandler) VerifUser(c echo.Context) error {
	var req binder.VerifUser
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	if err := h.userService.VerifUser(req.VerifCode); err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success verification user ", nil))
}
