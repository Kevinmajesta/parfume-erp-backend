package service

import (
	"errors"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
)

type VendorService interface {
	CreateVendor(vendor *entity.Vendors) (*entity.Vendors, error)
	CheckVendorExists(vendorId string) (bool, error)
	UpdateVendor(vendor *entity.Vendors) (*entity.Vendors, error)
	FindVendorByID(vendorId string) (*entity.Vendors, error)
	DeleteVendor(vendorId string) (bool, error)
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
	newVendor := entity.NewVendor(lastId, vendor.Vendorname, vendor.Addressone, vendor.Addresstwo, vendor.Phone, vendor.Email, vendor.Website)

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
