package service

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/email"
	"github.com/jung-kurt/gofpdf"
)

type RfqService interface {
	CreateRfq(rfq *entity.Rfqs) (*entity.Rfqs, error)
	GetCheckIDProduct(productId string) (bool, error)
	GetCheckIDMaterial(materialId string) (bool, error)
	UpdateRfq(rfq *entity.Rfqs) (*entity.Rfqs, error)
	FindRfqById(rfqId string) (*entity.Rfqs, error)
	UpdateRfqStatus(moId string) (*entity.Rfqs, error)
	FindAllRfq(page int) ([]entity.Rfqs, error)
	FindAllRfqBill(page int) ([]entity.Rfqs, error)
	CalculateOverview(rfqId string) (map[string]interface{}, error)
	GetEmailByVendorId(vendorId string) (string, error)
	SendRfqEmail(rfqId string, recipientEmail string) error
	DeleteRFQ(MoId string) (bool, error)
	CreateRfqPDF(rfqId string, recipientEmail string) ([]byte, error)
	UpdateRfqAll(rfqId string, updatedRfq *entity.Rfqs) (*entity.Rfqs, error)
	DeleteProductsByRfqId(rfqId string) error
}

type rfqService struct {
	rfqRepository  repository.RfqRepository
	rfqProductRepo repository.RfqProductRepository
	emailSender    *email.EmailSender
}

func NewRfqService(rfqRepository repository.RfqRepository, rfqProductRepo repository.RfqProductRepository, emailSender *email.EmailSender) *rfqService {
	return &rfqService{
		rfqRepository:  rfqRepository,
		rfqProductRepo: rfqProductRepo,
		emailSender:    emailSender,
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

func (s *rfqService) UpdateRfq(updatedRfq *entity.Rfqs) (*entity.Rfqs, error) {
	// Perbarui informasi RFQ
	existingRfq, err := s.rfqRepository.GetRfqById(updatedRfq.RfqId)
	if err != nil {
		return nil, err
	}
	if existingRfq == nil {
		return nil, fmt.Errorf("RFQ with id %s not found", updatedRfq.RfqId)
	}

	// Update detail RFQ
	existingRfq.OrderDate = updatedRfq.OrderDate
	existingRfq.Status = updatedRfq.Status
	existingRfq.VendorId = updatedRfq.VendorId
	existingRfq.UpdatedAt = time.Now()

	// Simpan perubahan RFQ
	updatedRfqResult, err := s.rfqRepository.UpdateRfq(existingRfq)
	if err != nil {
		return nil, err
	}

	// Hapus produk lama
	err = s.rfqProductRepo.DeleteProductsByRfqId(updatedRfq.RfqId)
	if err != nil {
		return nil, err
	}

	// Tambahkan produk baru
	for _, product := range updatedRfq.Products {
		lastProductId, err := s.rfqProductRepo.GetLastProductId()
		if err != nil {
			return nil, err
		}

		product.RfqsProductId = generateRfqsProductId(lastProductId)
		product.RfqId = updatedRfq.RfqId

		_, err = s.rfqProductRepo.CreateProduct(&product)
		if err != nil {
			return nil, err
		}
	}

	// Ambil produk terbaru untuk RFQ
	products, err := s.rfqProductRepo.GetProductsByRfqId(updatedRfq.RfqId)
	if err != nil {
		return nil, err
	}
	updatedRfqResult.Products = products

	return updatedRfqResult, nil
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

func (s *rfqService) UpdateRfqStatus(rfqId string) (*entity.Rfqs, error) {
	// Fetch the existing Manufacture Order
	mo, err := s.rfqRepository.GetRfqById(rfqId)
	if err != nil {
		return nil, errors.New("RFQ not found")
	}

	// Cycle through statuses
	switch mo.Status {
	case "RFQ":
		mo.Status = "Purchase Order"
	case "Purchase Order":
		mo.Status = "Recived"
	case "Recived":
		mo.Status = "Done"
	default:
		return nil, errors.New("invalid status transition")
	}

	// Save the updated Manufacture Order
	updatedMo, err := s.rfqRepository.UpdateRfqStatus(mo)
	if err != nil {
		return nil, errors.New("failed to update RFQ status")
	}

	return updatedMo, nil
}

func (s *rfqService) FindAllRfq(page int) ([]entity.Rfqs, error) {
	return s.rfqRepository.FindAllRfq(page)
}

func (s *rfqService) FindAllRfqBill(page int) ([]entity.Rfqs, error) {
	return s.rfqRepository.FindAllRfqBill(page)
}

func (s *rfqService) CalculateOverview(rfqId string) (map[string]interface{}, error) {
	// Ambil data RFQ dari repository
	rfq, err := s.rfqRepository.GetRfqById(rfqId)
	if err != nil {
		return nil, err
	}

	overview := make(map[string]interface{})
	overview["id_rfq"] = rfq.RfqId
	overview["vendor_id"] = rfq.VendorId
	overview["order_date"] = rfq.OrderDate
	overview["status"] = rfq.Status
	overview["created_at"] = rfq.CreatedAt
	overview["updated_at"] = rfq.UpdatedAt
	overview["deleted_at"] = rfq.DeletedAt
	overview["products"] = []map[string]interface{}{}
	totalCost := 0.0

	// Iterasi produk pada RFQ untuk menghitung detail
	for _, product := range rfq.Products {
		productDetails, err := s.rfqProductRepo.GetProductDetails(product.ProductId)
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

func (s *rfqService) GetEmailByVendorId(vendorId string) (string, error) {
	// Call the repository to check if email exists
	email, err := s.rfqRepository.CheckEmailExistsByVendorId(vendorId)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve email for vendor ID %s: %v", vendorId, err)
	}
	return email, nil
}

func (s *rfqService) SendRfqEmail(rfqId string, recipientEmail string) error {
	// Cari RFQ berdasarkan ID
	rfq, err := s.FindRfqById(rfqId)
	if err != nil {
		return fmt.Errorf("failed to find RFQ with id %s: %v", rfqId, err)
	}

	// Pastikan data produk tersedia
	if len(rfq.Products) == 0 {
		return errors.New("no products associated with this RFQ")
	}

	// Kirim email menggunakan service email (pastikan sudah diinisialisasi)
	err = s.emailSender.SendRfqEmail(
		recipientEmail, // Pass the email recipient as string
		rfq.RfqId,      // RFQ ID
		rfq.VendorId,   // Vendor ID
		rfq.OrderDate,  // Order Date
		rfq.Status,     // Status
		rfq.Products,   // List of products
	)
	if err != nil {
		return fmt.Errorf("failed to send RFQ email: %v", err)
	}

	return nil
}

func (s *rfqService) DeleteRFQ(MoId string) (bool, error) {
	material, err := s.rfqRepository.GetRfqById(MoId)
	if err != nil {
		return false, err
	}

	return s.rfqRepository.DeleteRfq(material)
}

func (s *rfqService) CreateRfqPDF(rfqId string, recipientEmail string) ([]byte, error) {
	// Cari RFQ berdasarkan ID
	rfq, err := s.FindRfqById(rfqId)
	if err != nil {
		return nil, fmt.Errorf("failed to find RFQ with id %s: %v", rfqId, err)
	}

	// Pastikan data produk tersedia
	if len(rfq.Products) == 0 {
		return nil, errors.New("no products associated with this RFQ")
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.AddPage()

	// Set Header
	pdf.SetFillColor(0, 102, 204)   // Header color
	pdf.SetTextColor(255, 255, 255) // White text
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 15, "RFQ Confirmation | RFQ ID: "+rfqId, "0", 1, "C", true, 0, "")
	pdf.Ln(10)

	// Add a horizontal line after the header
	pdf.SetDrawColor(0, 102, 204)
	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(10)

	// Reset font and text color for the body
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 12)

	// Greeting message
	pdf.Cell(0, 10, fmt.Sprintf("Dear Vendor,"))
	pdf.Ln(6)
	pdf.Cell(0, 10, "Here are the details of your RFQ:")
	pdf.Ln(10)

	// RFQ Table
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "RFQ ID:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, rfq.RfqId)
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Vendor ID:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, rfq.VendorId)
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Order Date:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, rfq.OrderDate)
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Status:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, rfq.Status)
	pdf.Ln(10)

	// Products Table
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Products:")
	pdf.Ln(6)

	// Table Header
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(50, 10, "Product Name")
	pdf.Cell(30, 10, "Quantity")
	pdf.Cell(40, 10, "Unit Price")
	pdf.Cell(40, 10, "Subtotal")
	pdf.Ln(6)

	// Table Body (Products)
	pdf.SetFont("Arial", "", 12)
	for _, product := range rfq.Products {
		pdf.Cell(50, 10, product.ProductName)
		pdf.Cell(30, 10, fmt.Sprintf("%s", product.Quantity))
		pdf.Cell(40, 10, fmt.Sprintf("%s", product.UnitPrice))
		pdf.Cell(40, 10, fmt.Sprintf("%s", product.Subtotal))
		pdf.Ln(6)
	}

	// Footer
	pdf.SetY(-30)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, "Thank you for doing business with us!")
	pdf.Ln(6)
	pdf.CellFormat(0, 10, "Best regards, Depublic Team", "", 0, "C", false, 0, "")

	// Save the PDF to a buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	// Return the PDF as a byte slice
	return buf.Bytes(), nil
}

func (s *rfqService) UpdateRfqAll(rfqId string, updatedRfq *entity.Rfqs) (*entity.Rfqs, error) {
	// Periksa apakah RFQ ada
	existingRfq, err := s.rfqRepository.GetRfqById(rfqId)
	if err != nil {
		return nil, err
	}
	if existingRfq == nil {
		return nil, fmt.Errorf("RFQ with id %s not found", rfqId)
	}

	// Update informasi RFQ
	existingRfq.OrderDate = updatedRfq.OrderDate
	existingRfq.Status = updatedRfq.Status
	existingRfq.VendorId = updatedRfq.VendorId
	existingRfq.UpdatedAt = time.Now()

	// Simpan perubahan RFQ
	savedRfq, err := s.rfqRepository.UpdateRfq(existingRfq)
	if err != nil {
		return nil, err
	}

	// Hapus produk lama jika ada produk baru yang diberikan
	if len(updatedRfq.Products) > 0 {
		err = s.rfqProductRepo.DeleteProductsByRfqId(rfqId)
		if err != nil {
			return nil, err
		}

		// Tambahkan produk baru
		for _, product := range updatedRfq.Products {
			lastProductId, err := s.rfqProductRepo.GetLastProductId()
			if err != nil {
				return nil, err
			}

			product.RfqsProductId = generateRfqsProductId(lastProductId)
			product.RfqId = rfqId

			_, err = s.rfqProductRepo.CreateProduct(&product)
			if err != nil {
				return nil, err
			}
		}
	}

	// Ambil produk terbaru untuk RFQ
	products, err := s.rfqProductRepo.GetProductsByRfqId(rfqId)
	if err != nil {
		return nil, err
	}
	savedRfq.Products = products

	return savedRfq, nil
}

func (s *rfqService) DeleteProductsByRfqId(rfqId string) error {
	return s.rfqProductRepo.DeleteProductsByRfqId(rfqId)
}
