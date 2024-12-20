package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/http/binder"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/service"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/response"
	"github.com/labstack/echo/v4"
)

type CostumerHandler struct {
	costumerService service.CostumerService
}

func NewCostumerHandler(costumerService service.CostumerService) CostumerHandler {
	return CostumerHandler{costumerService: costumerService}
}

func (h *CostumerHandler) CreateCostumer(c echo.Context) error {
	input := binder.CostumerCreateRequest{}

	// Bind input JSON ke struct
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "There is an input error"))
	}

	// Buat entitas produk dari input
	newVendor := &entity.Costumers{
		Costumername: input.Costumername,
		Addressone:   input.Addressone,
		Addresstwo:   input.Addresstwo,
		Phone:        input.Phone,
		Email:        input.Email,
		Status:       input.Status,
		Zip:          input.Zip,
		City:         input.City,
		Country:      input.Country,
		State:        input.State,
	}

	// Panggil service untuk membuat produk baru
	vendor, err := h.costumerService.CreateCostumer(newVendor)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Response jika berhasil
	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "Successfully input a new Costumer", vendor))
}

func (h *CostumerHandler) UpdateCostumer(c echo.Context) error {
	var input binder.CostumerUpdateRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "there is an input error"))
	}

	if input.CostumerId == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Costumer ID cannot be empty"))
	}
	exists, err := h.costumerService.CheckCostumerExists(input.CostumerId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "could not verify Costumer existence"))
	}
	if !exists {
		return c.JSON(http.StatusNotFound, response.ErrorResponse(http.StatusNotFound, "Costumer ID does not exist"))
	}

	inputUser := entity.UpdateCostumer(input.CostumerId, input.Costumername, input.Addressone, input.Addresstwo, input.Phone, input.Email, input.Status, input.State, input.Zip, input.Country, input.City)

	updatedProduk, err := h.costumerService.UpdateCostumer(inputUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success update Costumer", updatedProduk))
}
func (h *CostumerHandler) DeleteCostumer(c echo.Context) error {
	var input binder.CostumerDeleteRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "ada kesalahan input"))
	}

	isDeleted, err := h.costumerService.DeleteCostumer(input.CostumerId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "sukses delete costumer", isDeleted))
}
func (h *CostumerHandler) FindAllCostumer(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1 // default page if page parameter is invalid
	}

	materials, err := h.costumerService.FindAllCostumer(page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "success show data Costumers", materials))
}

func (h *CostumerHandler) GetCostumerProfile(c echo.Context) error {
	material_ID := c.Param("id_costumer")

	material, err := h.costumerService.FindCostumerBy(material_ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, "Failed to get costumer"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse(http.StatusOK, "successfully displays costumer data", material))
}

func (h *CostumerHandler) HandleCreateCostumerPDF(c echo.Context) error {
	// Get the Costumer ID from the URL parameters
	costumerId := c.Param("id_costumer")

	// Call the service method to get the costumer details
	costumer, err := h.costumerService.FindCostumerBy(costumerId) // This should return *entity.Costumers, not bool
	if err != nil {
		// If an error occurs while fetching the costumer, return error response
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Call the service method to create the PDF for the given costumer
	pdfBytes, err := h.costumerService.CreateCostumerPDF(costumer)
	if err != nil {
		// If an error occurs while generating the PDF, return error response
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Set the appropriate headers for PDF download
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=costumer-"+costumerId+".pdf")
	c.Response().Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes))) // Fix Content-Length header

	// Write the PDF content to the response
	if _, err := c.Response().Write(pdfBytes); err != nil {
		// Return error if there is an issue writing the PDF content
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return nil
}

func (h *CostumerHandler) HandleCreateCostumerPDFAll(c echo.Context) error {
	log.Println("Starting HandleCreateCostumerPDFAll")

	// Fetch all vendors from the service
	vendors, err := h.costumerService.FindAllCostumer(1)
	if err != nil {
		log.Println("FindAllCostumer error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Error fetching vendor list",
		})
	}

	log.Printf("Fetched %d vendors\n", len(vendors))

	// Convert vendors to []*entity.Costumers
	vendorPtrs := make([]*entity.Costumers, len(vendors))
	for i := range vendors {
		vendorPtrs[i] = &vendors[i]
	}

	// Generate the PDF for all vendors
	pdfData, err := h.costumerService.CreateCostumerPDFAll(vendorPtrs)
	if err != nil {
		log.Println("CreateCostumerPDFAll error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Error generating PDF",
		})
	}

	log.Printf("Generated PDF data of length: %d\n", len(pdfData))

	// Serve the PDF to the client
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=all_costumers.pdf")
	_, writeErr := c.Response().Write(pdfData)
	if writeErr != nil {
		log.Println("Error writing response:", writeErr)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Error writing PDF response",
		})
	}

	log.Println("PDF response successfully sent")
	return nil
}
