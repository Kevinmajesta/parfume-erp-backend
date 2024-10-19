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
	DeleteBom(bomId string) (bool, error)
	UpdateBOM(bom *entity.Bom) (*entity.Bom, error)
	CheckDuplicateProductInBOM(productId string, bomId string) (bool, error)
	CheckDuplicateMaterialInBOM(materialId string, bomId string) (bool, error)
	GetBOMByID(bomId string) (*entity.Bom, error)
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
