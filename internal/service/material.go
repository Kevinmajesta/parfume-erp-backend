package service

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code93"
	"github.com/jung-kurt/gofpdf"
)

type MaterialService interface {
	CreateMaterial(material *entity.Materials) (*entity.Materials, error)
	UpdateMaterial(material *entity.Materials) (*entity.Materials, error)
	CheckMaterialExists(materialId string) (bool, error)
	FindMaterialByID(materialId string) (*entity.Materials, error)
	DeleteMaterial(materialId string) (bool, error)
	FindAllMaterial(page int) ([]entity.Materials, error)
	SearchMaterialsByName(name string) ([]entity.Materials, error)
	GenerateAllMaterialsPDF(page int) (string, error)
	GenerateBarcodePDF(id string) (string, error)
	GenerateBarcode(id string) (string, error)
	GenerateMaterialPDF(id string) (string, error)
	GetMaterialByID(materialId string) (*entity.Materials, error)
	ReduceMaterialQty(input entity.Materials) (*entity.Materials, error)
	IncreaseMaterialQty(input entity.Materials) (*entity.Materials, error)
}

type materialService struct {
	materialRepository repository.MaterialRepository
}

func NewMaterialService(materialRepository repository.MaterialRepository) *materialService {
	return &materialService{
		materialRepository: materialRepository,
	}
}

func (s *materialService) CreateMaterial(material *entity.Materials) (*entity.Materials, error) {
	if material.Materialname == "" {
		return nil, errors.New("Materialname cannot be empty")
	}
	if material.Materialcategory == "" {
		return nil, errors.New("Materialcategory cannot be empty")
	}
	if material.Sellprice == "" {
		return nil, errors.New("Sellprice cannot be empty")
	}
	if material.Makeprice == "" {
		return nil, errors.New("Makeprice cannot be empty")
	}
	if material.Unit == "" {
		return nil, errors.New("Unit cannot be empty")
	}
	if material.Description == "" {
		return nil, errors.New("Description cannot be empty")
	}
	if material.Image == "" {
		return nil, errors.New("Image cannot be empty")
	}

	// Mendapatkan produk terakhir untuk generate ProdukId baru
	lastId, err := s.materialRepository.GetLastMaterial()
	if err != nil {
		return nil, err
	}

	// Buat produk baru dengan ID yang di-generate
	newMaterial := entity.NewMaterials(lastId, material.Materialname, material.Materialcategory, material.Sellprice, material.Makeprice, material.Unit, material.Description, material.Image)

	// Simpan produk ke database
	savedMaterial, err := s.materialRepository.CreateMaterial(newMaterial)
	if err != nil {
		return nil, err
	}

	return savedMaterial, nil
}

func (s *materialService) CheckMaterialExists(materialId string) (bool, error) {
	return s.materialRepository.CheckMaterialExists(materialId)
}

func (s *materialService) UpdateMaterial(material *entity.Materials) (*entity.Materials, error) {
	if material.Materialname == "" {
		return nil, errors.New("Materialname cannot be empty")
	}
	if material.Materialcategory == "" {
		return nil, errors.New("Materialcategory cannot be empty")
	}
	if material.Sellprice == "" {
		return nil, errors.New("Sellprice cannot be empty")
	}
	if material.Makeprice == "" {
		return nil, errors.New("Makeprice cannot be empty")
	}
	if material.Unit == "" {
		return nil, errors.New("Unit cannot be empty")
	}
	if material.Description == "" {
		return nil, errors.New("Description cannot be empty")
	}

	updatedMaterial, err := s.materialRepository.UpdateMaterial(material)
	if err != nil {
		return nil, err
	}

	return updatedMaterial, nil
}

func (s *materialService) FindMaterialByID(materialId string) (*entity.Materials, error) {
	return s.materialRepository.FindMaterialByID(materialId)
}

func (s *materialService) DeleteMaterial(materialId string) (bool, error) {
	material, err := s.materialRepository.FindMaterialByID(materialId)
	if err != nil {
		return false, err
	}

	log.Printf("Material to be deleted: %v", material)
	return s.materialRepository.DeleteMaterial(material)
}

func (s *materialService) FindAllMaterial(page int) ([]entity.Materials, error) {
	return s.materialRepository.FindAllMaterial(page)
}

func (s *materialService) SearchMaterialsByName(name string) ([]entity.Materials, error) {
	return s.materialRepository.SearchByName(name)
}

func (s *materialService) GenerateMaterialPDF(id string) (string, error) {
	material, err := s.materialRepository.FindMaterialByID(id)
	if err != nil {
		return "", err
	}

	// Konversi SellPrice dan MakePrice dari string ke float64
	sellPrice, err := strconv.ParseFloat(material.Sellprice, 64)
	if err != nil {
		return "", fmt.Errorf("invalid sell price: %v", err)
	}

	makePrice, err := strconv.ParseFloat(material.Makeprice, 64)
	if err != nil {
		return "", fmt.Errorf("invalid make price: %v", err)
	}

	// Buat barcode terlebih dahulu
	barcodeFile, err := s.GenerateBarcode(id)
	if err != nil {
		return "", err
	}

	// Membuat PDF dengan layout yang lebih estetik
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.AddPage()

	// Header dengan latar belakang berwarna
	pdf.SetFillColor(100, 150, 255) // Warna biru muda
	pdf.SetTextColor(255, 255, 255) // Warna putih
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 12, "Material Details", "0", 1, "C", true, 0, "")
	pdf.Ln(10)

	// Reset warna teks ke hitam
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 12)

	// Label dan Nilai untuk setiap field
	labels := []string{"Material ID:", "Material Name:", "Category:", "Sell Price:", "Make Price:", "Unit:", "Description:"}
	values := []string{
		material.MaterialId,
		material.Materialname,
		material.Materialcategory,
		fmt.Sprintf("%.2f", sellPrice),
		fmt.Sprintf("%.2f", makePrice),
		material.Unit,
		material.Description,
	}

	for i := 0; i < len(labels); i++ {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 8, labels[i])
		pdf.SetFont("Arial", "", 12)
		pdf.MultiCell(0, 8, values[i], "", "", false)
	}

	// Menambahkan barcode ke PDF
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Material Barcode")
	pdf.Ln(15)
	pdf.Image(barcodeFile, 15, pdf.GetY(), 60, 20, false, "", 0, "")

	// Simpan PDF ke file
	fileName := "material_" + material.MaterialId + ".pdf"
	err = pdf.OutputFileAndClose(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *materialService) GenerateBarcode(id string) (string, error) {
	material, err := s.materialRepository.FindMaterialByID(id)
	if err != nil {
		return "", err
	}

	// Generate the barcode image using code93 encoding
	qrCode, _ := code93.Encode(material.MaterialId, true, true)

	// Ensure the barcode image is scaled to a proper size
	qrCode, _ = barcode.Scale(qrCode, 650, 250)

	// Convert the image to an 8-bit-per-channel image format compatible with PNG
	// Create a new RGBA image from the original QR code
	barcodeImage := image.NewRGBA(qrCode.Bounds())

	// Copy the QR code pixels into the new image (if necessary)
	for y := 0; y < qrCode.Bounds().Dy(); y++ {
		for x := 0; x < qrCode.Bounds().Dx(); x++ {
			barcodeImage.Set(x, y, qrCode.At(x, y))
		}
	}

	// Generate the file path for the barcode
	barcodeFile := "barcode_" + material.MaterialId + ".png"

	// Create the PNG file
	file, err := os.Create(barcodeFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Encode the image as a PNG
	err = png.Encode(file, barcodeImage)
	if err != nil {
		return "", err
	}

	return barcodeFile, nil
}

func (s *materialService) GenerateBarcodePDF(id string) (string, error) {
	material, err := s.materialRepository.FindMaterialByID(id)
	if err != nil {
		return "", err
	}

	// Buat barcode terlebih dahulu
	barcodeFile, err := s.GenerateBarcode(id)
	if err != nil {
		return "", err
	}

	// Membuat PDF untuk barcode
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Product Barcode")

	pdf.Ln(12)
	pdf.Image(barcodeFile, 10, 20, 40, 40, false, "", 0, "")

	barcodePDF := "barcode_" + material.MaterialId + ".pdf"
	err = pdf.OutputFileAndClose(barcodePDF)
	if err != nil {
		return "", err
	}

	return barcodePDF, nil
}

func (s *materialService) GenerateAllMaterialsPDF(page int) (string, error) {
	materials, err := s.FindAllMaterial(page)
	if err != nil {
		return "", err
	}

	// Create a new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.AddPage()

	// Header with background color
	pdf.SetFillColor(100, 150, 255) // Light blue color
	pdf.SetTextColor(255, 255, 255) // White text
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 12, "All Materials List - Page "+fmt.Sprintf("%d", page), "0", 1, "C", true, 0, "")
	pdf.Ln(10)

	// Reset text color to black
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 12)

	// Loop through each material and add its details to the PDF
	for _, material := range materials {
		// Convert SellPrice and MakePrice from string to float64
		sellPrice, err := strconv.ParseFloat(material.Sellprice, 64)
		if err != nil {
			return "", errors.New("invalid sell price for material: " + material.Materialname)
		}

		makePrice, err := strconv.ParseFloat(material.Makeprice, 64)
		if err != nil {
			return "", errors.New("invalid make price for material: " + material.Materialname)
		}

		// Generate barcode for the material
		barcodeFile, err := s.GenerateBarcode(material.MaterialId)
		if err != nil {
			return "", fmt.Errorf("failed to generate barcode for material: %v", material.MaterialId)
		}

		// Labels and values for each material
		labels := []string{"Material ID:", "Material Name:", "Category:", "Sell Price:", "Make Price:", "Unit:", "Description:"}
		values := []string{
			material.MaterialId,
			material.Materialname,
			material.Materialcategory,
			fmt.Sprintf("%.2f", sellPrice),
			fmt.Sprintf("%.2f", makePrice),
			material.Unit,
			material.Description,
		}

		// Add each label and value to the PDF
		for i := 0; i < len(labels); i++ {
			pdf.SetFont("Arial", "B", 12)
			pdf.Cell(40, 8, labels[i])
			pdf.SetFont("Arial", "", 12)
			pdf.MultiCell(0, 8, values[i], "", "", false)
		}

		// Add the barcode image to the PDF
		pdf.Ln(3)
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, "Material Barcode")
		pdf.Ln(15)
		pdf.Image(barcodeFile, 15, pdf.GetY(), 60, 20, false, "", 0, "")

		// Add a separator line after each material
		pdf.Ln(10)
		pdf.Cell(0,29, "----------------------------------------")
		pdf.Ln(20)
	}

	// Save the PDF to a file
	fileName := "all_materials_page_" + fmt.Sprintf("%d", page) + ".pdf"
	err = pdf.OutputFileAndClose(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *materialService) GetMaterialByID(materialId string) (*entity.Materials, error) {

	material, err := s.materialRepository.FindMaterialByID(materialId)
	if err != nil {
		return nil, err
	}
	return material, nil
}

func (s *materialService) ReduceMaterialQty(input entity.Materials) (*entity.Materials, error) {
	// Fetch material by ID
	material, err := s.materialRepository.FindMaterialByName(input.MaterialId)
	if err != nil {
		return nil, errors.New("material not found")
	}

	// Check if there's enough quantity
	if material.Qty < input.Qty {
		return nil, errors.New("insufficient quantity")
	}

	// Reduce the quantity
	material.Qty -= input.Qty // This works as both are float64

	// Update material in database
	err = s.materialRepository.Update(material)
	if err != nil {
		return nil, errors.New("failed to update material")
	}

	return material, nil
}

func (s *materialService) IncreaseMaterialQty(input entity.Materials) (*entity.Materials, error) {
	// Fetch material by ID
	material, err := s.materialRepository.FindMaterialByID(input.MaterialId)
	if err != nil {
		return nil, errors.New("material not found")
	}

	// Reduce the quantity
	material.Qty += input.Qty // This works as both are float64

	// Update material in database
	err = s.materialRepository.Update(material)
	if err != nil {
		return nil, errors.New("failed to update material")
	}

	return material, nil
}
