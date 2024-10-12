package service

import (
	"errors"
	"log"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
)

type ProductService interface {
	CreateProduct(product *entity.Products) (*entity.Products, error)
	UpdateProduct(product *entity.Products) (*entity.Products, error)
	CheckProductExists(productId string) (bool, error)
	FindProductByID(productId string) (*entity.Products, error)
	DeleteProduct(productId string) (bool, error)
	FindAllProduct(page int) ([]entity.Products, error)
	SearchProductsByName(name string) ([]entity.Products, error)
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