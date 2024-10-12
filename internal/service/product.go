package service

import (
	"errors"
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
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
	GenerateProductPDF(id string) (string, error)
	GenerateBarcodePDF(id string) (string, error)
	GenerateBarcode(id string) (string, error)
	GenerateAllProductsPDF(page int) (string, error)
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

	// Mendapatkan produk terakhir untuk generate ProdukId baru
	lastId, err := s.productRepository.GetLastProduct()
	if err != nil {
		return nil, err
	}

	// Buat produk baru dengan ID yang di-generate
	newProduct := entity.NewProduct(lastId, product.Productname, product.Productcategory, product.Sellprice, product.Makeprice, product.Pajak, product.Description, product.Image)

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
	if product.Image == "" {
		return nil, errors.New("Image cannot be empty")
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

func (s *productService) SearchProductsByName(name string) ([]entity.Products, error) {
	return s.productRepository.SearchByName(name)
}

func (s *productService) GenerateProductPDF(id string) (string, error) {
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

	// Membuat PDF dengan gofpdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Product Details")

	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, "Product ID: "+product.ProdukId)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Product Name: "+product.Productname)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Category: "+product.Productcategory)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Sell Price: "+fmt.Sprintf("%.2f", sellPrice))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Make Price: "+fmt.Sprintf("%.2f", makePrice))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Tax: "+fmt.Sprintf("%.2f%%", pajak))
	pdf.Ln(10)
	pdf.Cell(40, 10, "Description: "+product.Description)

	// Simpan PDF ke file
	fileName := "product_" + product.ProdukId + ".pdf"
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

	qrCode, _ := qr.Encode(product.ProdukId, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 200, 200) // Resize barcode

	barcodeFile := "barcode_" + product.ProdukId + ".png"
	file, _ := os.Create(barcodeFile)
	defer file.Close()

	png.Encode(file, qrCode)

	return barcodeFile, nil
}

// GenerateBarcodePDF untuk memasukkan barcode ke dalam PDF
func (s *productService) GenerateBarcodePDF(id string) (string, error) {
	product, err := s.productRepository.FindProductByID(id)
	if err != nil {
		return "", err
	}

	// Buat barcode terlebih dahulu
	barcodeFile, err := s.GenerateBarcode(id)
	if err != nil {
		return "", err
	}

	// Membuat PDF untuk barcode
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Product Barcode")

	pdf.Ln(12)
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

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "All Products List")

	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)

	for _, product := range products {
		// Konversi SellPrice, MakePrice, dan Pajak dari string ke float64
		sellPrice, err := strconv.ParseFloat(product.Sellprice, 64)
		if err != nil {
			return "", errors.New("invalid sell price for product: " + product.Productname)
		}

		makePrice, err := strconv.ParseFloat(product.Makeprice, 64)
		if err != nil {
			return "", errors.New("invalid make price for product: " + product.Productname)
		}

		pajak, err := strconv.ParseFloat(product.Pajak, 64)
		if err != nil {
			return "", errors.New("invalid tax for product: " + product.Productname)
		}

		// Menambahkan informasi produk ke PDF
		pdf.Cell(40, 10, "Product ID: "+product.ProdukId)
		pdf.Ln(10)
		pdf.Cell(40, 10, "Product Name: "+product.Productname)
		pdf.Ln(10)
		pdf.Cell(40, 10, "Category: "+product.Productcategory)
		pdf.Ln(10)
		pdf.Cell(40, 10, "Sell Price: "+fmt.Sprintf("%.2f", sellPrice))
		pdf.Ln(10)
		pdf.Cell(40, 10, "Make Price: "+fmt.Sprintf("%.2f", makePrice))
		pdf.Ln(10)
		pdf.Cell(40, 10, "Tax: "+fmt.Sprintf("%.2f", pajak))
		pdf.Ln(10)
		pdf.Cell(40, 10, "Description: "+product.Description)
		pdf.Ln(10)
		pdf.Cell(40, 10, "-------------------------------------")
		pdf.Ln(10)
	}

	fileName := "all_products_page_" + fmt.Sprintf("%d", page) + ".pdf"
	err = pdf.OutputFileAndClose(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
