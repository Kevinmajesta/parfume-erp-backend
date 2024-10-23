package service

import (
	"errors"
	"fmt"
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

	// Konversi SellPrice, MakePrice, dan Pajak dari string ke float64
	sellPrice, err := strconv.ParseFloat(material.Sellprice, 64)
	if err != nil {
		return "", fmt.Errorf("invalid sell price: %v", err)
	}

	makePrice, err := strconv.ParseFloat(material.Makeprice, 64)
	if err != nil {
		return "", fmt.Errorf("invalid make price: %v", err)
	}

	// Membuat PDF dengan gofpdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Product Details")

	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Material ID: "+material.MaterialId)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Material Name: "+material.Materialname)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Category: "+material.Materialcategory)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Sell Price: "+fmt.Sprintf("%.2f", sellPrice))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Make Price: "+fmt.Sprintf("%.2f", makePrice))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Unit: "+material.Unit)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Description: "+material.Description)

	// Simpan PDF ke file
	fileName := "product_" + material.MaterialId + ".pdf"
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

	qrCode, _ := code93.Encode(material.MaterialId, true, true)
	qrCode, _ = barcode.Scale(qrCode, 650, 250)

	barcodeFile := "barcode_" + material.MaterialId + ".png"
	file, _ := os.Create(barcodeFile)
	defer file.Close()

	png.Encode(file, qrCode)

	return barcodeFile, nil
}

// GenerateBarcodePDF untuk memasukkan barcode ke dalam PDF
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
	material, err := s.FindAllMaterial(page)
	if err != nil {
		return "", err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "All Material List")

	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)

	for _, material := range material {
		// Konversi SellPrice, MakePrice, dan Pajak dari string ke float64
		sellPrice, err := strconv.ParseFloat(material.Sellprice, 64)
		if err != nil {
			return "", errors.New("invalid sell price for material: " + material.Materialname)
		}

		makePrice, err := strconv.ParseFloat(material.Makeprice, 64)
		if err != nil {
			return "", errors.New("invalid make price for material: " + material.Materialname)
		}

		// Menambahkan informasi produk ke PDF
		pdf.Cell(40, 10, "Material ID: "+material.MaterialId)
		pdf.Ln(10)
		pdf.Cell(40, 10, "Material Name: "+material.Materialname)
		pdf.Ln(10)
		pdf.Cell(40, 10, "Category: "+material.Materialname)
		pdf.Ln(10)
		pdf.Cell(40, 10, "Sell Price: "+fmt.Sprintf("%.2f", sellPrice))
		pdf.Ln(10)
		pdf.Cell(40, 10, "Make Price: "+fmt.Sprintf("%.2f", makePrice))
		pdf.Ln(10)
		pdf.Cell(40, 10, "Unit: "+material.Unit)
		pdf.Ln(10)
		pdf.Cell(40, 10, "Description: "+material.Description)
		pdf.Ln(10)
		pdf.Cell(40, 10, "-------------------------------------")
		pdf.Ln(10)
	}

	fileName := "all_products_page_" + fmt.Sprintf("%d", page) + ".pdf"
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
