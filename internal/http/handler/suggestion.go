package handler

import (
	"net/http"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/service"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/response"
	"github.com/labstack/echo/v4"
)

type SuggestionHandler struct {
	suggestionService service.SuggestionService
	userService       service.UserService
}

func NewSuggestionHandler(suggestionService service.SuggestionService, userService service.UserService) SuggestionHandler {
	return SuggestionHandler{suggestionService: suggestionService, userService: userService}
}

func (h *SuggestionHandler) CreateSuggestion(c echo.Context) error {
	var input entity.Suggestion

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid request"))
	}

	// Set the type and message for the notification, assuming they come from the input
	notification := &entity.Suggestion{
		Type:    input.Type,
		Message: input.Message,
	}

	if input.Type == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Suggestion type cannot be empty"))
	}
	if input.Message == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Suggestion message cannot be empty"))
	}

	// Call the service to create the notification for all users
	if err := h.suggestionService.CreateSuggestion(notification); err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusCreated, "Suggestion created successfully", nil))
}
