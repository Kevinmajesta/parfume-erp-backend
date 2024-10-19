package service

import (
	"fmt"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
)

type BOMService interface {
	CreateBOM(bom *entity.Bom) (*entity.Bom, error)
	GetCheckIDProduct(productId string) (bool, error)
	GetCheckIDMaterial(materialId string) (bool, error)
	FindAllBom(page int) ([]entity.Bom, error)
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
