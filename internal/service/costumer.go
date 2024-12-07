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
