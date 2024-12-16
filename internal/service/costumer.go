package service

import (
	"errors"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
)

type CostumerService interface {
	CreateCostumer(vendor *entity.Costumers) (*entity.Costumers, error)
	CheckCostumerExists(vendorId string) (bool, error)
	UpdateCostumer(vendor *entity.Costumers) (*entity.Costumers, error)
	DeleteCostumer(vendorId string) (bool, error)
	FindAllCostumer(page int) ([]entity.Costumers, error)
	FindCostumerBy(materialId string) (*entity.Costumers, error)
	// GetCostumerById(vendorId string) (bool, error)
	// CreateCostumerPDF(costumer *entity.Costumers) ([]byte, error)
}

type costumerService struct {
	costumerRepository repository.CostumerRepository
}

func NewCostumerService(costumerRepository repository.CostumerRepository) *costumerService {
	return &costumerService{
		costumerRepository: costumerRepository,
	}
}

func (s *costumerService) CreateCostumer(vendor *entity.Costumers) (*entity.Costumers, error) {
	if vendor.Costumername == "" {
		return nil, errors.New("Vendorname cannot be empty")
	}
	if vendor.Addressone == "" {
		return nil, errors.New("ProduAddressonectcategory cannot be empty")
	}
	if vendor.Addresstwo == "" {
		return nil, errors.New("Addresstwo cannot be empty")
	}
	if vendor.Phone == "" {
		return nil, errors.New("Phone cannot be empty")
	}
	if vendor.Email == "" {
		return nil, errors.New("Email cannot be empty")
	}

	lastId, err := s.costumerRepository.GetLastCostumer()
	if err != nil {
		return nil, err
	}

	// Buat produk baru dengan ID yang di-generate
	newVendor := entity.NewCostumer(lastId, vendor.Costumername, vendor.Addressone, vendor.Addresstwo, vendor.Phone, vendor.Email, vendor.Status, vendor.Zip, vendor.City, vendor.Country, vendor.State)

	// Simpan produk ke database
	savedVendor, err := s.costumerRepository.CreateCostumer(newVendor)
	if err != nil {
		return nil, err
	}

	return savedVendor, nil
}

func (s *costumerService) CheckCostumerExists(vendorId string) (bool, error) {
	return s.costumerRepository.CheckCostumerExists(vendorId)
}

func (s *costumerService) UpdateCostumer(vendor *entity.Costumers) (*entity.Costumers, error) {
	if vendor.Costumername == "" {
		return nil, errors.New("Costumername cannot be empty")
	}
	if vendor.Addressone == "" {
		return nil, errors.New("Addressone cannot be empty")
	}
	if vendor.Addresstwo == "" {
		return nil, errors.New("Addresstwo cannot be empty")
	}
	if vendor.Phone == "" {
		return nil, errors.New("Phone cannot be empty")
	}
	if vendor.Email == "" {
		return nil, errors.New("Email cannot be empty")
	}

	updatedVendor, err := s.costumerRepository.UpdateCostumer(vendor)
	if err != nil {
		return nil, err
	}

	return updatedVendor, nil
}

func (s *costumerService) DeleteCostumer(vendorId string) (bool, error) {
	vendor, err := s.costumerRepository.FindCostumerByID(vendorId)
	if err != nil {
		return false, err
	}

	return s.costumerRepository.DeleteCostumer(vendor)
}

func (s *costumerService) FindAllCostumer(page int) ([]entity.Costumers, error) {
	return s.costumerRepository.FindAllCostumer(page)
}

func (s *costumerService) FindCostumerBy(materialId string) (*entity.Costumers, error) {
	return s.costumerRepository.FindCostumerByID(materialId)
}

// func (s *costumerService) CreateCostumerPDF(costumer *entity.Costumers) ([]byte, error) {
// 	// Validate costumer details
// 	if costumer.Costumername == "" {
// 		return nil, errors.New("Costumer name cannot be empty")
// 	}

// 	// Create a new PDF document
// 	pdf := gofpdf.New("P", "mm", "A4", "")
// 	pdf.SetMargins(15, 20, 15)
// 	pdf.AddPage()

// 	// Set Header
// 	pdf.SetFillColor(100, 150, 255) // Light blue color
// 	pdf.SetTextColor(255, 255, 255) // White text
// 	pdf.SetFont("Arial", "B", 24)
// 	pdf.CellFormat(0, 12, "Costumer Details", "0", 1, "C", true, 0, "")
// 	pdf.Ln(10)

// 	// Add a horizontal line
// 	pdf.SetDrawColor(100, 150, 255)
// 	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
// 	pdf.Ln(5)

// 	// Reset font and text color for body
// 	pdf.SetTextColor(0, 0, 0)
// 	pdf.SetFont("Arial", "", 12)

// 	// Customer Details
// 	labels := []string{
// 		"Costumer ID:", "Costumer Name:", "Address Line 1:", "Address Line 2:",
// 		"Phone:", "Email:", "Status:", "Zip:", "City:", "Country:", "State:"}

// 	values := []string{
// 		costumer.CostumerId, costumer.Costumername, costumer.Addressone, costumer.Addresstwo,
// 		costumer.Phone, costumer.Email, costumer.Status, costumer.Zip,
// 		costumer.City, costumer.Country, costumer.State}

// 	// Loop through labels and values
// 	for i := 0; i < len(labels); i++ {
// 		// Label in bold font
// 		pdf.SetFont("Arial", "B", 12)
// 		pdf.Cell(40, 10, labels[i])
// 		// Value in normal font
// 		pdf.SetFont("Arial", "", 12)
// 		pdf.MultiCell(0, 10, values[i], "", "", false)
// 		pdf.Ln(4) // Line break after each field
// 	}

// 	// Add a footer
// 	pdf.SetY(-30) // Footer at the bottom
// 	pdf.SetFont("Arial", "I", 10)
// 	pdf.Cell(0, 10, "Generated by Your Company Name")

// 	// Save the PDF to a buffer
// 	var buf bytes.Buffer
// 	err := pdf.Output(&buf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Return the PDF as a byte slice
// 	return buf.Bytes(), nil
// }
