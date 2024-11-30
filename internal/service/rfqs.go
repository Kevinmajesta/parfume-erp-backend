package service

import (
	"fmt"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
)

type RfqService interface {
	CreateRfq(rfq *entity.Rfqs) (*entity.Rfqs, error)
	GetCheckIDProduct(productId string) (bool, error)
	GetCheckIDMaterial(materialId string) (bool, error)
	UpdateRfq(rfq *entity.Rfqs) (*entity.Rfqs, error)
	FindRfqById(rfqId string) (*entity.Rfqs, error)
}

type rfqService struct {
	rfqRepository  repository.RfqRepository
	rfqProductRepo repository.RfqProductRepository
}

func NewRfqService(rfqRepository repository.RfqRepository, rfqProductRepo repository.RfqProductRepository) *rfqService {
	return &rfqService{
		rfqRepository:  rfqRepository,
		rfqProductRepo: rfqProductRepo,
	}
}

func generateRfqsId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "RFQ-%d", &newNumber)
		newNumber++
	}

	return fmt.Sprintf("RFQ-%05d", newNumber)
}

func generateRfqsProductId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "R-%d", &newNumber)
		newNumber++
	}

	return fmt.Sprintf("R-%05d", newNumber)
}

func (s *rfqService) CreateRfq(rfq *entity.Rfqs) (*entity.Rfqs, error) {

	lastId, err := s.rfqRepository.GetLastRfq()
	if err != nil {
		return nil, err
	}

	newRfq := entity.NewRfqs(lastId, rfq.OrderDate, rfq.Status, rfq.VendorId)

	savedRfq, err := s.rfqRepository.CreateRfq(newRfq)
	if err != nil {
		return nil, err
	}

	if len(rfq.Products) > 0 {
		for _, products := range rfq.Products {

			lastProductId, err := s.rfqProductRepo.GetLastProductId()
			if err != nil {
				return nil, err
			}

			products.RfqsProductId = generateRfqsProductId(lastProductId)
			products.RfqId = savedRfq.RfqId

			_, err = s.rfqProductRepo.CreateProduct(&products)
			if err != nil {
				return nil, err
			}
		}
	}

	products, err := s.rfqProductRepo.GetProductsByRfqId(savedRfq.RfqId)
	if err != nil {
		return nil, err
	}
	savedRfq.Products = products

	return savedRfq, nil
}

func (s *rfqService) GetCheckIDProduct(productId string) (bool, error) {
	return s.rfqRepository.CheckProductExists(productId)
}

func (s *rfqService) GetCheckIDMaterial(materialId string) (bool, error) {
	return s.rfqProductRepo.CheckMaterialExists(materialId)
}

func (s *rfqService) UpdateRfq(rfq *entity.Rfqs) (*entity.Rfqs, error) {
	// Ambil data RFQ lama
	existingRfq, err := s.rfqRepository.GetRfqById(rfq.RfqId)
	if err != nil {
		return nil, err
	}

	if existingRfq == nil {
		return nil, fmt.Errorf("RFQ with id %s not found", rfq.RfqId)
	}

	// Update hanya data RFQ tanpa menyentuh produk
	updatedRfq := entity.UpdateRfqs(
		rfq.RfqId,
		rfq.OrderDate,
		rfq.Status,
		rfq.VendorId,
		existingRfq.Status,
	)

	// Perbarui data RFQ
	result, err := s.rfqRepository.UpdateRfq(updatedRfq)
	if err != nil {
		return nil, err
	}

	// Update produk yang sudah ada
	// Update produk yang sudah ada
	// Update produk yang sudah ada
	for _, product := range rfq.Products {
		// Periksa apakah produk sudah ada berdasarkan RfqId dan ProductId
		existingProduct, err := s.rfqProductRepo.GetProductByRfqIdAndProductId(rfq.RfqId, product.ProductId)
		if err != nil {
			return nil, err
		}

		if existingProduct != nil {
			// Jika produk sudah ada, lakukan pembaruan
			product.RfqsProductId = existingProduct.RfqsProductId
			product.RfqId = rfq.RfqId

			_, err = s.rfqProductRepo.UpdateProduct(&product)
			if err != nil {
				return nil, err
			}
		} else {
			// Jika produk belum ada, buat produk baru
			lastProductId, err := s.rfqProductRepo.GetLastProductId()
			if err != nil {
				return nil, err
			}

			product.RfqsProductId = generateRfqsProductId(lastProductId)
			product.RfqId = rfq.RfqId

			_, err = s.rfqProductRepo.CreateProduct(&product)
			if err != nil {
				return nil, err
			}
		}
	}

	// Ambil produk yang diperbarui
	updatedProducts, err := s.rfqProductRepo.GetProductsByRfqId(rfq.RfqId)
	if err != nil {
		return nil, err
	}
	result.Products = updatedProducts

	return result, nil
}

func (s *rfqService) FindRfqById(rfqId string) (*entity.Rfqs, error) {
	// Ambil data RFQ berdasarkan ID
	rfq, err := s.rfqRepository.GetRfqById(rfqId)
	if err != nil {
		return nil, err
	}

	// Jika RFQ tidak ditemukan, return nil
	if rfq == nil {
		return nil, fmt.Errorf("RFQ with id %s not found", rfqId)
	}

	// Ambil produk yang terkait dengan RFQ
	products, err := s.rfqProductRepo.GetProductsByRfqId(rfqId)
	if err != nil {
		return nil, err
	}

	// Tambahkan produk ke RFQ
	rfq.Products = products

	return rfq, nil
}
