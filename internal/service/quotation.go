package service

import (
	"errors"
	"fmt"
	"strconv"

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
	DeleteRFQ(MoId string) (bool, error)
	FindAllQuo(page int) ([]entity.Quotations, error)
	UpdateQuoStatus(rfqId string) (*entity.Quotations, error)
	CalculateOverview(rfqId string) (map[string]interface{}, error)
	FindAllQuoBill(page int) ([]entity.Quotations, error)
	SendQuoEmail(rfqId string, recipientEmail string) error
	GetEmailByCostumerId(vendorId string) (string, error)
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

func (s *quoService) DeleteRFQ(MoId string) (bool, error) {
	material, err := s.quoRepository.GetQuoById(MoId)
	if err != nil {
		return false, err
	}

	return s.quoRepository.DeleteQuo(material)
}

func (s *quoService) FindAllQuo(page int) ([]entity.Quotations, error) {
	return s.quoRepository.FindAllQuo(page)
}

func (s *quoService) FindAllQuoBill(page int) ([]entity.Quotations, error) {
	return s.quoRepository.FindAllQuoBill(page)
}

func (s *quoService) UpdateQuoStatus(rfqId string) (*entity.Quotations, error) {
	// Fetch the existing Manufacture Order
	mo, err := s.quoRepository.GetQuoById(rfqId)
	if err != nil {
		return nil, errors.New("Quo not found")
	}

	// Cycle through statuses
	switch mo.Status {
	case "QUOTATION":
		mo.Status = "Sales Order"
	case "Sales Order":
		mo.Status = "Invoiced"
	case "Invoiced":
		mo.Status = "Delivery"
	case "Delivery":
		mo.Status = "Done"
	default:
		return nil, errors.New("invalid status transition")
	}

	// Save the updated Manufacture Order
	updatedMo, err := s.quoRepository.UpdateQuoStatus(mo)
	if err != nil {
		return nil, errors.New("failed to update Quo status")
	}

	return updatedMo, nil
}

func (s *quoService) CalculateOverview(rfqId string) (map[string]interface{}, error) {
	// Ambil data RFQ dari repository
	rfq, err := s.quoRepository.GetQuoById(rfqId)
	if err != nil {
		return nil, err
	}

	overview := make(map[string]interface{})
	overview["id_quotation"] = rfq.QuotationsId
	overview["costumer_id"] = rfq.CostumerId
	overview["order_date"] = rfq.OrderDate
	overview["status"] = rfq.Status
	overview["created_at"] = rfq.CreatedAt
	overview["updated_at"] = rfq.UpdatedAt
	overview["deleted_at"] = rfq.DeletedAt
	overview["products"] = []map[string]interface{}{}
	totalCost := 0.0

	// Iterasi produk pada RFQ untuk menghitung detail
	for _, product := range rfq.Products {
		productDetails, err := s.quoProductRepo.GetProductDetails(product.ProductId)
		if err != nil {
			return nil, err
		}

		// Kalkulasi harga subtotal
		subtotal, err := strconv.ParseFloat(product.Subtotal, 64)
		if err != nil {
			return nil, err
		}

		quantity, err := strconv.ParseFloat(product.Quantity, 64)
		if err != nil {
			return nil, err
		}
		productCost := subtotal * quantity

		productDetail := map[string]interface{}{
			"product_id":   product.ProductId,
			"product_name": product.ProductName,
			"quantity":     product.Quantity,
			"unit_price":   product.UnitPrice,
			"tax":          product.Tax,
			"subtotal":     product.Subtotal,
			"total_cost":   fmt.Sprintf("Rp %.2f", productCost),
			"vendor_price": productDetails.Sellprice, // Jika ingin menggunakan harga dari vendor
		}

		// Hitung total biaya
		totalCost += productCost

		overview["products"] = append(overview["products"].([]map[string]interface{}), productDetail)
	}

	overview["total_cost"] = fmt.Sprintf("Rp %.2f", totalCost)
	return overview, nil
}

func (s *quoService) SendQuoEmail(rfqId string, recipientEmail string) error {
	// Cari RFQ berdasarkan ID
	rfq, err := s.FindQuoById(rfqId)
	if err != nil {
		return fmt.Errorf("failed to find Quo with id %s: %v", rfqId, err)
	}

	// Pastikan data produk tersedia
	if len(rfq.Products) == 0 {
		return errors.New("no products associated with this RFQ")
	}

	// Kirim email menggunakan service email (pastikan sudah diinisialisasi)
	err = s.emailSender.SendQuoEmail(
		recipientEmail,   // Pass the email recipient as string
		rfq.QuotationsId, // RFQ ID
		rfq.CostumerId,   // Vendor ID
		rfq.OrderDate,    // Order Date
		rfq.Status,       // Status
		rfq.Products,     // List of products
	)
	if err != nil {
		return fmt.Errorf("failed to send RFQ email: %v", err)
	}

	return nil
}

func (s *quoService) GetEmailByCostumerId(vendorId string) (string, error) {
	// Call the repository to check if email exists
	email, err := s.quoRepository.CheckEmailExistsByCostumerId(vendorId)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve email for Costumer ID %s: %v", vendorId, err)
	}
	return email, nil
}
