package service

import (
	"bytes"
	"fmt"
	"strings"

	"strconv"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
	"github.com/jung-kurt/gofpdf"
)

type BOMService interface {
	CreateBOM(bom *entity.Bom) (*entity.Bom, error)
	GetCheckIDProduct(productId string) (bool, error)
	GetCheckIDMaterial(materialId string) (bool, error)
	FindAllBom(page int) ([]entity.Bom, error)
	DeleteBom(bomId string) (bool, error)
	UpdateBOM(bom *entity.Bom) (*entity.Bom, error)
	CheckDuplicateProductInBOM(productId string, bomId string) (bool, error)
	CheckDuplicateMaterialInBOM(materialId string, bomId string) (bool, error)
	GetBOMByID(bomId string) (*entity.Bom, error)
	CalculateOverview(bomId string) (map[string]interface{}, error)
	GenerateBOMPDF(overview map[string]interface{}) ([]byte, error)
}

type bomService struct {
	bomRepo         repository.BOMRepository
	bomMaterialRepo repository.BOMMaterialRepository
}

func NewBOMService(bomRepo repository.BOMRepository, bomMaterialRepo repository.BOMMaterialRepository) BOMService {
	return &bomService{
		bomRepo:         bomRepo,
		bomMaterialRepo: bomMaterialRepo,
	}
}

func generateBOMMaterialId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "BMM-%d", &newNumber)
		newNumber++
	}
	return fmt.Sprintf("BMM-%05d", newNumber)
}

func (s *bomService) CreateBOM(bom *entity.Bom) (*entity.Bom, error) {

	lastId, err := s.bomRepo.GetLastBomId()
	if err != nil {
		return nil, err
	}

	newBom := entity.NewBom(lastId, bom.ProductId, bom.ProductName, bom.ProductPreference, bom.Quantity)

	savedBom, err := s.bomRepo.CreateBOM(newBom)
	if err != nil {
		return nil, err
	}

	if len(bom.Materials) > 0 {
		for _, material := range bom.Materials {

			lastMaterialId, err := s.bomMaterialRepo.GetLastMaterialId()
			if err != nil {
				return nil, err
			}

			material.IdBomMaterial = generateBOMMaterialId(lastMaterialId)
			material.BomId = savedBom.BomId

			_, err = s.bomMaterialRepo.CreateMaterial(&material)
			if err != nil {
				return nil, err
			}
		}
	}

	materials, err := s.bomMaterialRepo.GetMaterialsByBomId(savedBom.BomId)
	if err != nil {
		return nil, err
	}
	savedBom.Materials = materials

	return savedBom, nil
}

func (s *bomService) GetCheckIDProduct(productId string) (bool, error) {
	return s.bomRepo.CheckProductExists(productId)
}

func (s *bomService) GetCheckIDMaterial(materialId string) (bool, error) {
	return s.bomMaterialRepo.CheckMaterialExists(materialId)
}

func (s *bomService) FindAllBom(page int) ([]entity.Bom, error) {
	return s.bomRepo.FindAllBom(page)
}

func (s *bomService) DeleteBom(bomId string) (bool, error) {
	return s.bomRepo.DeleteBom(bomId)
}

func (s *bomService) UpdateBOM(bom *entity.Bom) (*entity.Bom, error) {
	// Update BOM in database
	updatedBom, err := s.bomRepo.UpdateBOM(bom)
	if err != nil {
		return nil, err
	}

	// Handle materials update
	if len(bom.Materials) > 0 {
		// First, delete the old materials and add new ones
		if err := s.bomMaterialRepo.DeleteMaterialsByBomId(bom.BomId); err != nil {
			return nil, err
		}

		// Insert updated materials
		for _, material := range bom.Materials {
			lastMaterialId, err := s.bomMaterialRepo.GetLastMaterialId()
			if err != nil {
				return nil, err
			}

			material.IdBomMaterial = generateBOMMaterialId(lastMaterialId)
			material.BomId = bom.BomId

			_, err = s.bomMaterialRepo.CreateMaterial(&material)
			if err != nil {
				return nil, err
			}
		}
	}

	// Fetch the updated materials to return along with the updated BOM
	materials, err := s.bomMaterialRepo.GetMaterialsByBomId(bom.BomId)
	if err != nil {
		return nil, err
	}
	updatedBom.Materials = materials

	return updatedBom, nil
}

func (s *bomService) CheckDuplicateProductInBOM(productId string, bomId string) (bool, error) {
	existingBom, err := s.bomRepo.FindBOMByProductIDAndBOMID(productId, bomId)
	if err != nil {
		return false, err
	}

	if existingBom != nil {
		return true, nil
	}

	return false, nil
}

func (s *bomService) CheckDuplicateMaterialInBOM(materialId string, bomId string) (bool, error) {
	existingBom, err := s.bomMaterialRepo.FindBOMByMaterialIDAndBOMID(materialId, bomId)
	if err != nil {
		return false, err
	}

	if existingBom != nil {
		return true, nil
	}

	return false, nil
}

// Service Layer
func (s *bomService) GetBOMByID(bomId string) (*entity.Bom, error) {
	return s.bomRepo.FindBOMByID(bomId)
}

func (s *bomService) CalculateOverview(bomId string) (map[string]interface{}, error) {
	bom, err := s.bomRepo.FindBOMByID(bomId)
	if err != nil {
		return nil, err
	}

	overview := make(map[string]interface{})
	overview["bom_id"] = bom.BomId
	overview["product_name"] = bom.ProductName
	overview["materials"] = []map[string]interface{}{}
	totalCost := 0.0

	// Fetch product details once to use for all materials
	productDetails, err := s.bomRepo.GetProductDetails(bom.ProductId)
	if err != nil {
		return nil, err
	}

	// Prepare product details for overview
	overview["product_details"] = map[string]interface{}{
		"make_price": productDetails.Makeprice,
		"sell_price": productDetails.Sellprice,
	}

	for _, material := range bom.Materials {
		// Fetch material price
		materialDetails, err := s.bomMaterialRepo.GetMaterialDetails(material.IdMaterial)
		if err != nil {
			return nil, err
		}

		// Prepare material detail for overview
		quantity, err := strconv.ParseFloat(material.Quantity, 64)
		if err != nil {
			return nil, err
		}
		makePrice, err := strconv.ParseFloat(materialDetails.Makeprice, 64)
		if err != nil {
			return nil, err
		}

		// Calculate product cost based on quantity
		productCost := makePrice * quantity

		materialDetail := map[string]interface{}{
			"material":     material.MaterialName,
			"quantity":     material.Quantity,
			"product_cost": fmt.Sprintf("%.2f", productCost), // Calculated value
			"bom_cost":     materialDetails.Sellprice,        // Use sell price or whatever is appropriate
		}

		// Calculate total cost
		totalCost += productCost

		overview["materials"] = append(overview["materials"].([]map[string]interface{}), materialDetail)
	}

	overview["total_cost"] = fmt.Sprintf("Rp %.2f", totalCost)
	return overview, nil
}

func (s *bomService) GenerateBOMPDF(overview map[string]interface{}) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.AddPage()

	// Set a simple header
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 12, "BOM Overview", "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Reset text color and font
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 12)

	// Add BOM ID and product name
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "BOM ID:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, overview["bom_id"].(string))
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Product Name:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, overview["product_name"].(string))
	pdf.Ln(12)

	// Add product details section
	productDetails := overview["product_details"].(map[string]interface{})
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Make Price:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, productDetails["make_price"].(string))
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Sell Price:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, productDetails["sell_price"].(string))
	pdf.Ln(12)

	// Materials Table Header with subtle styling
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(70, 10, "Material")
	pdf.Cell(30, 10, "Quantity")
	pdf.Cell(50, 10, "Product Cost")
	pdf.Cell(50, 10, "BOM Cost")
	pdf.Ln(10)

	// Materials Table Rows with basic formatting
	pdf.SetFont("Arial", "", 12)
	materials := overview["materials"].([]map[string]interface{})
	for _, material := range materials {
		pdf.Cell(70, 10, material["material"].(string))
		pdf.Cell(30, 10, material["quantity"].(string))
		pdf.Cell(50, 10, material["product_cost"].(string))
		pdf.Cell(50, 10, material["bom_cost"].(string))
		pdf.Ln(10)
	}

	// Add a line separator after materials
	pdf.Ln(5)
	pdf.Cell(0, 1, "")
	pdf.Ln(10)

	// Total Cost Section
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Total Cost:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, overview["total_cost"].(string))
	pdf.Ln(15)

	// Output PDF to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Fungsi untuk parsing biaya
func parseCost(costStr string) (float64, error) {
	// Hapus "Rp " dan konversi menjadi float64
	cleanedCost := strings.ReplaceAll(costStr, "Rp ", "")
	cleanedCost = strings.ReplaceAll(cleanedCost, ".", "") // Hapus titik jika ada
	return strconv.ParseFloat(cleanedCost, 64)
}
