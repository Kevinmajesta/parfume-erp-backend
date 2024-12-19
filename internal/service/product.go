package service

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code93"
	"github.com/jung-kurt/gofpdf"
)

type ProductService interface {
	CreateProduct(product *entity.Products) (*entity.Products, error)
	UpdateProduct(product *entity.Products) (*entity.Products, error)
	CheckProductExists(productId string) (bool, error)
	FindProductByID(productId string) (*entity.Products, error)
	DeleteProduct(productId string) (bool, error)
	FindAllProduct(page int) ([]entity.Products, error)
	SearchProductsByName(name string) ([]entity.Products, error)
	GenerateProductPDFWithBarcode(id string) (string, error)
	GenerateBarcodePDF(id string) (string, error)
	GenerateBarcode(id string) (string, error)
	GenerateAllProductsPDF(page int) (string, error)
	FindAllProductVariant(page int) ([]entity.Products, error)
	GetProductByID(productId string) (*entity.Products, error)
	IncreaseProductQty(input entity.Products) (*entity.Products, error)
	DecreaseProductQty(input entity.Products) (*entity.Products, error)
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) *productService {
	return &productService{
		productRepository: productRepository,
	}
}

func (s *productService) CreateProduct(product *entity.Products) (*entity.Products, error) {
	if product.Productname == "" {
		return nil, errors.New("Productname cannot be empty")
	}
	if product.Productcategory == "" {
		return nil, errors.New("Productcategory cannot be empty")
	}
	if product.Sellprice == "" {
		return nil, errors.New("Sellprice cannot be empty")
	}
	if product.Makeprice == "" {
		return nil, errors.New("Makeprice cannot be empty")
	}
	if product.Pajak == "" {
		return nil, errors.New("Pajak cannot be empty")
	}
	if product.Description == "" {
		return nil, errors.New("Description cannot be empty")
	}
	if product.Image == "" {
		return nil, errors.New("Image cannot be empty")
	}

	if product.Variant != "yes" && product.Variant != "no" {
		return nil, errors.New("Variant must be either 'yes' or 'no'")
	}

	// Mendapatkan produk terakhir untuk generate ProdukId baru
	lastId, err := s.productRepository.GetLastProduct()
	if err != nil {
		return nil, err
	}

	// Buat produk baru dengan ID yang di-generate
	newProduct := entity.NewProduct(lastId, product.Productname, product.Productcategory, product.Sellprice, product.Makeprice, product.Pajak, product.Description, product.Image, product.Variant)

	// Simpan produk ke database
	savedProduct, err := s.productRepository.CreateProduct(newProduct)
	if err != nil {
		return nil, err
	}

	return savedProduct, nil
}

func (s *productService) CheckProductExists(productId string) (bool, error) {
	return s.productRepository.CheckProductExists(productId)
}

func (s *productService) UpdateProduct(product *entity.Products) (*entity.Products, error) {
	if product.Productname == "" {
		return nil, errors.New("Productname cannot be empty")
	}
	if product.Productcategory == "" {
		return nil, errors.New("Productcategory cannot be empty")
	}
	if product.Sellprice == "" {
		return nil, errors.New("Sellprice cannot be empty")
	}
	if product.Makeprice == "" {
		return nil, errors.New("Makeprice cannot be empty")
	}
	if product.Pajak == "" {
		return nil, errors.New("Pajak cannot be empty")
	}
	if product.Description == "" {
		return nil, errors.New("Description cannot be empty")
	}

	if product.Variant != "yes" && product.Variant != "no" {
		return nil, errors.New("Variant must be either 'yes' or 'no'")
	}

	updatedProduct, err := s.productRepository.UpdateProduct(product)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s *productService) FindProductByID(productId string) (*entity.Products, error) {
	return s.productRepository.FindProductByID(productId)
}

func (s *productService) DeleteProduct(productId string) (bool, error) {
	product, err := s.productRepository.FindProductByID(productId)
	if err != nil {
		return false, err
	}

	log.Printf("Product to be deleted: %v", product)
	return s.productRepository.DeleteProduct(product)
}

func (s *productService) FindAllProduct(page int) ([]entity.Products, error) {
	return s.productRepository.FindAllProduct(page)
}

func (s *productService) FindAllProductVariant(page int) ([]entity.Products, error) {
	return s.productRepository.FindAllProductVariant(page)
}

func (s *productService) SearchProductsByName(name string) ([]entity.Products, error) {
	return s.productRepository.SearchByName(name)
}

func (s *productService) GenerateProductPDFWithBarcode(id string) (string, error) {
	product, err := s.productRepository.FindProductByID(id)
	if err != nil {
		return "", err
	}

	// Konversi SellPrice, MakePrice, dan Pajak dari string ke float64
	sellPrice, err := strconv.ParseFloat(product.Sellprice, 64)
	if err != nil {
		return "", fmt.Errorf("invalid sell price: %v", err)
	}

	makePrice, err := strconv.ParseFloat(product.Makeprice, 64)
	if err != nil {
		return "", fmt.Errorf("invalid make price: %v", err)
	}

	pajak, err := strconv.ParseFloat(product.Pajak, 64)
	if err != nil {
		return "", fmt.Errorf("invalid tax: %v", err)
	}

	// Generate Barcode
	barcodeFile, err := s.GenerateBarcode(id)
	if err != nil {
		return "", err
	}

	// Membuat PDF dengan gofpdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.AddPage()

	// Header dengan background warna
	pdf.SetFillColor(100, 150, 255) // Warna biru muda
	pdf.SetTextColor(255, 255, 255) // Warna putih
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 12, "Product Details", "0", 1, "C", true, 0, "")
	pdf.Ln(10)

	// Reset warna teks ke hitam
	pdf.SetTextColor(0, 0, 0)

	// Informasi produk
	pdf.SetFont("Arial", "", 12)
	infoLabels := []string{"Product ID:", "Product Name:", "Category:", "Sell Price:", "Make Price:", "Tax:", "Description:"}
	infoValues := []string{
		product.ProdukId,
		product.Productname,
		product.Productcategory,
		fmt.Sprintf("%.2f", sellPrice),
		fmt.Sprintf("%.2f", makePrice),
		fmt.Sprintf("%.2f%%", pajak),
		product.Description,
	}

	for i := 0; i < len(infoLabels); i++ {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(40, 10, infoLabels[i])
		pdf.SetFont("Arial", "", 12)
		pdf.MultiCell(0, 10, infoValues[i], "", "", false)
	}

	// Garis pemisah
	pdf.Ln(5)
	pdf.SetDrawColor(200, 200, 200) // Warna abu-abu terang
	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(10)

	// Section Barcode
	pdf.SetFillColor(230, 230, 230) // Background abu-abu terang
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "Product Barcode", "0", 1, "C", true, 0, "")
	pdf.Ln(10)

	// Menambahkan Barcode
	pdf.Image(barcodeFile, 55, pdf.GetY(), 100, 40, false, "", 0, "")
	pdf.Ln(50)

	// Simpan PDF ke file
	fileName := "product_with_barcode_" + product.ProdukId + ".pdf"
	err = pdf.OutputFileAndClose(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *productService) GenerateBarcode(id string) (string, error) {
	product, err := s.productRepository.FindProductByID(id)
	if err != nil {
		return "", err
	}

	// Generate barcode
	qrCode, err := code93.Encode(product.ProdukId, true, true)
	if err != nil {
		return "", fmt.Errorf("failed to generate barcode: %v", err)
	}

	// Scale barcode
	qrCode, err = barcode.Scale(qrCode, 650, 250)
	if err != nil {
		return "", fmt.Errorf("failed to scale barcode: %v", err)
	}

	// Convert to 8-bit RGBA
	rgbaImage := image.NewRGBA(qrCode.Bounds())
	draw.Draw(rgbaImage, qrCode.Bounds(), qrCode, image.Point{}, draw.Src)

	// Save as PNG
	barcodeFile := "barcode_" + product.ProdukId + ".png"
	file, err := os.Create(barcodeFile)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Encode barcode as 8-bit PNG
	err = png.Encode(file, rgbaImage)
	if err != nil {
		return "", fmt.Errorf("failed to encode PNG: %v", err)
	}

	return barcodeFile, nil
}

func (s *productService) GenerateBarcodePDF(id string) (string, error) {
	product, err := s.productRepository.FindProductByID(id)
	if err != nil {
		return "", err
	}

	barcodeFile, err := s.GenerateBarcode(id)
	if err != nil {
		return "", err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Product Barcode")

	pdf.Ln(30)
	pdf.Image(barcodeFile, 10, 20, 40, 40, false, "", 0, "")

	barcodePDF := "barcode_" + product.ProdukId + ".pdf"
	err = pdf.OutputFileAndClose(barcodePDF)
	if err != nil {
		return "", err
	}

	return barcodePDF, nil
}

func (s *productService) GenerateAllProductsPDF(page int) (string, error) {
	products, err := s.FindAllProduct(page)
	if err != nil {
		return "", err
	}

	// Membuat PDF dengan layout yang lebih estetik
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.AddPage()

	// Header dengan latar belakang berwarna
	pdf.SetFillColor(100, 150, 255) // Warna biru muda
	pdf.SetTextColor(255, 255, 255) // Warna putih
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 12, "All Products List", "0", 1, "C", true, 0, "")
	pdf.Ln(10)

	// Reset warna teks ke hitam
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 12)

	// Menambahkan informasi produk
	for _, product := range products {
		// Konversi SellPrice, MakePrice, dan Pajak dari string ke float64
		sellPrice, err := strconv.ParseFloat(product.Sellprice, 64)
		if err != nil {
			return "", fmt.Errorf("invalid sell price for product: %s", product.Productname)
		}

		makePrice, err := strconv.ParseFloat(product.Makeprice, 64)
		if err != nil {
			return "", fmt.Errorf("invalid make price for product: %s", product.Productname)
		}

		pajak, err := strconv.ParseFloat(product.Pajak, 64)
		if err != nil {
			return "", fmt.Errorf("invalid tax for product: %s", product.Productname)
		}

		// Label dan Nilai untuk setiap field
		labels := []string{"Product ID:", "Product Name:", "Category:", "Sell Price:", "Make Price:", "Tax:", "Description:"}
		values := []string{
			product.ProdukId,
			product.Productname,
			product.Productcategory,
			fmt.Sprintf("%.2f", sellPrice),
			fmt.Sprintf("%.2f", makePrice),
			fmt.Sprintf("%.2f%%", pajak),
			product.Description,
		}

		for i := 0; i < len(labels); i++ {
			pdf.SetFont("Arial", "B", 12)
			pdf.Cell(40, 8, labels[i])
			pdf.SetFont("Arial", "", 12)
			pdf.MultiCell(0, 8, values[i], "", "", false)
		}

		// Generate barcode untuk produk ini
		barcodeFile, err := s.GenerateBarcode(product.ProdukId)
		if err != nil {
			return "", fmt.Errorf("failed to generate barcode for product: %s", product.Productname)
		}

		// Menambahkan barcode ke PDF
		pdf.Ln(5)
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 10, "Barcode:")
		pdf.Ln(12)
		pdf.Image(barcodeFile, 15, pdf.GetY(), 60, 20, false, "", 0, "")
		pdf.Ln(25)

		// Garis pemisah antar produk
		pdf.SetDrawColor(200, 200, 200)
		pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
		pdf.Ln(10)
	}

	// Nama file berdasarkan halaman
	fileName := fmt.Sprintf("all_products_page_%d.pdf", page)
	err = pdf.OutputFileAndClose(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *productService) GetProductByID(productId string) (*entity.Products, error) {

	product, err := s.productRepository.FindProductByID(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) IncreaseProductQty(input entity.Products) (*entity.Products, error) {
	// Fetch product by ID
	product, err := s.productRepository.FindProductByID(input.ProdukId)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Increase the quantity
	product.Qty += input.Qty // Ensure Qty is float64 in Product struct

	// Check if the new quantity is valid
	if product.Qty < 0 {
		return nil, errors.New("resulting quantity cannot be negative")
	}

	// Update product in the database
	if err := s.productRepository.Update(product); err != nil {
		return nil, errors.New("failed to update product quantity")
	}

	return product, nil
}

func (s *productService) DecreaseProductQty(input entity.Products) (*entity.Products, error) {
	// Fetch product by ID
	product, err := s.productRepository.FindProductByID(input.ProdukId)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Increase the quantity
	product.Qty -= input.Qty // Ensure Qty is float64 in Product struct

	// Update product in the database
	if err := s.productRepository.Update(product); err != nil {
		return nil, errors.New("failed to update product quantity")
	}

	return product, nil
}
