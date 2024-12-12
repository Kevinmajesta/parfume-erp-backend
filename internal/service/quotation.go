package service

import (
	"fmt"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/email"
)

type QuoService interface {
	CreateQuo(rfq *entity.Quotations) (*entity.Quotations, error)
	GetCheckIDProduct(productId string) (bool, error)
	GetCheckIDMaterial(materialId string) (bool, error)
	UpdateQuo(rfq *entity.Quotations) (*entity.Quotations, error)
	FindQuoById(rfqId string) (*entity.Quotations, error)
}

type quoService struct {
	quoRepository  repository.QuoRepository
	quoProductRepo repository.QuoProductRepository
	emailSender    *email.EmailSender
}

func NewQuoService(quoRepository repository.QuoRepository, quoProductRepo repository.QuoProductRepository, emailSender *email.EmailSender) *quoService {
	return &quoService{
		quoRepository:  quoRepository,
		quoProductRepo: quoProductRepo,
		emailSender:    emailSender,
	}
}

func generateQuosId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "QUO-%d", &newNumber)
		newNumber++
	}

	return fmt.Sprintf("QUO-%05d", newNumber)
}

func generateQuosProductId(lastId string) string {
	var newNumber int
	if lastId == "" {
		newNumber = 1
	} else {
		fmt.Sscanf(lastId, "QPR-%d", &newNumber)
		newNumber++
	}

	return fmt.Sprintf("QPR-%05d", newNumber)
}

func (s *quoService) CreateQuo(rfq *entity.Quotations) (*entity.Quotations, error) {

	lastId, err := s.quoRepository.GetLastQuo()
	if err != nil {
		return nil, err
	}

	newRfq := entity.NewQuo(lastId, rfq.OrderDate, rfq.Status, rfq.CostumerId)

	savedRfq, err := s.quoRepository.CreateQuo(newRfq)
	if err != nil {
		return nil, err
	}

	if len(rfq.Products) > 0 {
		for _, products := range rfq.Products {

			lastProductId, err := s.quoProductRepo.GetLastProductId()
			if err != nil {
				return nil, err
			}

			products.QuotationsProductId = generateQuosProductId(lastProductId)
			products.QuotationsId = savedRfq.QuotationsId

			_, err = s.quoProductRepo.CreateProduct(&products)
			if err != nil {
				return nil, err
			}
		}
	}

	products, err := s.quoProductRepo.GetProductsByQuoId(savedRfq.QuotationsId)
	if err != nil {
		return nil, err
	}
	savedRfq.Products = products

	return savedRfq, nil
}

func (s *quoService) GetCheckIDProduct(productId string) (bool, error) {
	return s.quoRepository.CheckProductExists(productId)
}

func (s *quoService) GetCheckIDMaterial(materialId string) (bool, error) {
	return s.quoProductRepo.CheckMaterialExists(materialId)
}

func (s *quoService) UpdateQuo(rfq *entity.Quotations) (*entity.Quotations, error) {
	// Ambil data RFQ lama
	existingRfq, err := s.quoRepository.GetQuoById(rfq.QuotationsId)
	if err != nil {
		return nil, err
	}

	if existingRfq == nil {
		return nil, fmt.Errorf("Quo with id %s not found", rfq.QuotationsId)
	}

	// Update hanya data RFQ tanpa menyentuh produk
	updatedRfq := entity.UpdateQuo(
		rfq.QuotationsId,
		rfq.OrderDate,
		rfq.Status,
		rfq.CostumerId,
		existingRfq.Status,
	)

	// Perbarui data RFQ
	result, err := s.quoRepository.UpdateQuo(updatedRfq)
	if err != nil {
		return nil, err
	}

	for _, product := range rfq.Products {
		// Periksa apakah produk sudah ada berdasarkan RfqId dan ProductId
		existingProduct, err := s.quoProductRepo.GetProductByQuoIdAndProductId(rfq.QuotationsId, product.ProductId)
		if err != nil {
			return nil, err
		}

		if existingProduct != nil {
			// Jika produk sudah ada, lakukan pembaruan
			product.QuotationsProductId = existingProduct.QuotationsProductId
			product.QuotationsId = rfq.QuotationsId
			product.CostumerId = rfq.CostumerId

			_, err = s.quoProductRepo.UpdateProduct(&product)
			if err != nil {
				return nil, err
			}
		} else {
			// Jika produk belum ada, buat produk baru
			lastProductId, err := s.quoProductRepo.GetLastProductId()
			if err != nil {
				return nil, err
			}

			product.QuotationsProductId = generateQuosProductId(lastProductId)
			product.QuotationsId = rfq.QuotationsId

			_, err = s.quoProductRepo.CreateProduct(&product)
			if err != nil {
				return nil, err
			}
		}
	}

	// Ambil produk yang diperbarui
	updatedProducts, err := s.quoProductRepo.GetProductsByQuoId(rfq.QuotationsId)
	if err != nil {
		return nil, err
	}
	result.Products = updatedProducts

	return result, nil
}

func (s *quoService) FindQuoById(rfqId string) (*entity.Quotations, error) {
	// Ambil data RFQ berdasarkan ID
	rfq, err := s.quoRepository.GetQuoById(rfqId)
	if err != nil {
		return nil, err
	}

	// Jika RFQ tidak ditemukan, return nil
	if rfq == nil {
		return nil, fmt.Errorf("Quo with id %s not found", rfqId)
	}

	// Ambil produk yang terkait dengan RFQ
	products, err := s.quoProductRepo.GetProductsByQuoId(rfqId)
	if err != nil {
		return nil, err
	}

	// Tambahkan produk ke RFQ
	rfq.Products = products

	return rfq, nil
}
