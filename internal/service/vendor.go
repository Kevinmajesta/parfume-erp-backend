package service

import (
	"bytes"
	"errors"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
	"github.com/jung-kurt/gofpdf"
)

type VendorService interface {
	CreateVendor(vendor *entity.Vendors) (*entity.Vendors, error)
	CheckVendorExists(vendorId string) (bool, error)
	UpdateVendor(vendor *entity.Vendors) (*entity.Vendors, error)
	FindVendorByID(vendorId string) (*entity.Vendors, error)
	DeleteVendor(vendorId string) (bool, error)
	FindVendorBy(materialId string) (*entity.Vendors, error)
	FindAllVendor(page int) ([]entity.Vendors, error)
	CreateVendorPDF(vendor *entity.Vendors) ([]byte, error)
}

type vendorService struct {
	vendorRepository repository.VendorRepository
}

func NewVendorService(vendorRepository repository.VendorRepository) *vendorService {
	return &vendorService{
		vendorRepository: vendorRepository,
	}
}

func (s *vendorService) CreateVendor(vendor *entity.Vendors) (*entity.Vendors, error) {
	if vendor.Vendorname == "" {
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
	if vendor.Website == "" {
		return nil, errors.New("Website cannot be empty")
	}

	lastId, err := s.vendorRepository.GetLastVendor()
	if err != nil {
		return nil, err
	}

	// Buat produk baru dengan ID yang di-generate
	newVendor := entity.NewVendor(lastId, vendor.Vendorname, vendor.Addressone, vendor.Addresstwo, vendor.Phone, vendor.Email, vendor.Website, vendor.Status, vendor.Zip, vendor.City, vendor.Country, vendor.State)

	// Simpan produk ke database
	savedVendor, err := s.vendorRepository.CreateVendor(newVendor)
	if err != nil {
		return nil, err
	}

	return savedVendor, nil
}

func (s *vendorService) CheckVendorExists(vendorId string) (bool, error) {
	return s.vendorRepository.CheckVendorExists(vendorId)
}

func (s *vendorService) UpdateVendor(vendor *entity.Vendors) (*entity.Vendors, error) {
	if vendor.Vendorname == "" {
		return nil, errors.New("Vendorname cannot be empty")
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
	if vendor.Website == "" {
		return nil, errors.New("Website cannot be empty")
	}

	updatedVendor, err := s.vendorRepository.UpdateVendor(vendor)
	if err != nil {
		return nil, err
	}

	return updatedVendor, nil
}

func (s *vendorService) FindVendorByID(vendorId string) (*entity.Vendors, error) {
	return s.vendorRepository.FindVendorByID(vendorId)
}

func (s *vendorService) DeleteVendor(vendorId string) (bool, error) {
	vendor, err := s.vendorRepository.FindVendorByID(vendorId)
	if err != nil {
		return false, err
	}

	return s.vendorRepository.DeleteVendor(vendor)
}

func (s *vendorService) FindAllVendor(page int) ([]entity.Vendors, error) {
	return s.vendorRepository.FindAllVendor(page)
}

func (s *vendorService) FindVendorBy(materialId string) (*entity.Vendors, error) {
	return s.vendorRepository.FindVendorByID(materialId)
}

func (s *vendorService) CreateVendorPDF(vendor *entity.Vendors) ([]byte, error) {
	// Check if vendor details are valid
	if vendor.Vendorname == "" {
		return nil, errors.New("Vendor name cannot be empty")
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.AddPage()

	// Set Header
	pdf.SetFillColor(0, 102, 204) // Light blue color for the header
	pdf.SetTextColor(255, 255, 255) // White text for the header
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 15, "Vendor Information", "0", 1, "C", true, 0, "")
	pdf.Ln(10)

	// Add a horizontal line after the header
	pdf.SetDrawColor(0, 102, 204)
	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(10)

	// Reset font and text color for the body
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 12)

	// Section Title - Vendor Details
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Vendor Details")
	pdf.Ln(6) // Line break after title

	// Create a table of details (Labels + Values)
	labels := []string{
		"Vendor ID:", "Vendor Name:", "Address Line 1:", "Address Line 2:",
		"Phone:", "Email:", "Website:", "Status:", "Zip:", "City:", "Country:", "State:"}

	values := []string{
		vendor.VendorId, vendor.Vendorname, vendor.Addressone, vendor.Addresstwo,
		vendor.Phone, vendor.Email, vendor.Website, vendor.Status,
		vendor.Zip, vendor.City, vendor.Country, vendor.State}

	// Adjust for label and value column widths
	labelWidth := 45.0
	valueWidth := 140.0

	for i := 0; i < len(labels); i++ {
		// Label in bold font
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(labelWidth, 10, labels[i])

		// Value in normal font, wrapping the text if too long
		pdf.SetFont("Arial", "", 12)
		pdf.MultiCell(valueWidth, 10, values[i], "", "", false)

		// Add a line break after each entry
		pdf.Ln(4)
	}

	// Add a footer (with company details and alignment)
	pdf.SetY(-30) // Position the footer at the bottom
	pdf.SetFont("Arial", "I", 10)
	pdf.SetTextColor(0, 102, 204) // Footer color
	pdf.Cell(0, 10, "Konate Parfume")

	// Save the PDF to a buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	// Return the PDF as a byte slice
	return buf.Bytes(), nil
}
