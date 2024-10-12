package service

import (
	"errors"
	"log"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
)

type MaterialService interface {
	CreateMaterial(material *entity.Materials) (*entity.Materials, error)
	UpdateMaterial(material *entity.Materials) (*entity.Materials, error)
	CheckMaterialExists(materialId string) (bool, error)
	FindMaterialByID(materialId string) (*entity.Materials, error)
	DeleteMaterial(materialId string) (bool, error)
	FindAllMaterial(page int) ([]entity.Materials, error)
	SearchMaterialsByName(name string) ([]entity.Materials, error)
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
	if material.Image == "" {
		return nil, errors.New("Image cannot be empty")
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
