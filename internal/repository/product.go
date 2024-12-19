package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/cache"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *entity.Products) (*entity.Products, error)
	GetLastProduct() (string, error)
	CheckProductExists(productId string) (bool, error)
	UpdateProduct(product *entity.Products) (*entity.Products, error)
	FindProductByID(productId string) (*entity.Products, error)
	DeleteProduct(product *entity.Products) (bool, error)
	FindAllProduct(page int) ([]entity.Products, error)
	SearchByName(name string) ([]entity.Products, error)
	FindAllProductVariant(page int) ([]entity.Products, error)
	Update(product *entity.Products) error
}

type productRepository struct {
	db        *gorm.DB
	cacheable cache.Cacheable
}

func NewProductRepository(db *gorm.DB, cacheable cache.Cacheable) *productRepository {
	return &productRepository{db: db, cacheable: cacheable}
}

func (r *productRepository) CreateProduct(product *entity.Products) (*entity.Products, error) {
	if err := r.db.Create(&product).Error; err != nil {
		return product, err
	}
	r.cacheable.Delete("FindAllProducts_page_1")
	r.cacheable.Delete("FindAllProducts_page_2")
	return product, nil
}

func (r *productRepository) GetLastProduct() (string, error) {
	var lastProduct entity.Products
	err := r.db.Order("id_product DESC").First(&lastProduct).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "PRF-00000", nil
	} else if err != nil {
		return "", err
	}

	return lastProduct.ProdukId, nil
}

func (r *productRepository) CheckProductExists(productId string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Products{}).Where("id_product = ?", productId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *productRepository) UpdateProduct(product *entity.Products) (*entity.Products, error) {
	fields := make(map[string]interface{})

	if product.Productname != "" {
		fields["productname"] = product.Productname
	}
	if product.Productcategory != "" {
		fields["productcategory"] = product.Productcategory
	}
	if product.Sellprice != "" {
		fields["sellprice"] = product.Sellprice
	}
	if product.Makeprice != "" {
		fields["makeprice"] = product.Makeprice
	}
	if product.Pajak != "" {
		fields["pajak"] = product.Pajak
	}
	if product.Description != "" {
		fields["description"] = product.Description
	}

	if err := r.db.Model(product).Where("id_product = ?", product.ProdukId).Updates(fields).Error; err != nil {
		return product, err
	}
	r.cacheable.Delete("FindAllProducts_page_1")

	return product, nil
}

func (r *productRepository) FindProductByID(productId string) (*entity.Products, error) {
	product := new(entity.Products)
	if err := r.db.Where("id_product = ?", productId).First(product).Error; err != nil {
		log.Printf("Error finding product by ID: %v", err)
		return nil, err // Pastikan mengembalikan nil, err
	}
	log.Printf("Product found: %v", product)
	return product, nil
}

func (r *productRepository) DeleteProduct(product *entity.Products) (bool, error) {
	log.Printf("Deleting product: %v", product)
	// Ensure hard delete by using Unscoped()
	if err := r.db.Unscoped().Delete(product).Error; err != nil {
		log.Printf("Error deleting product: %v", err)
		return false, err
	}
	log.Println("Product deleted successfully")
	r.cacheable.Delete("FindAllProducts_page_1")
	return true, nil
}

func (r *productRepository) FindAllProduct(page int) ([]entity.Products, error) {
	var products []entity.Products
	key := fmt.Sprintf("FindAllProducts_page_%d", page)
	const pageSize = 100

	data, _ := r.cacheable.Get(key)
	if data == "" {
		offset := (page - 1) * pageSize
		if err := r.db.Limit(pageSize).Offset(offset).Where("variant = ?", "no").Find(&products).Error; err != nil {
			return products, err
		}
		marshalledproducts, _ := json.Marshal(products)
		err := r.cacheable.Set(key, marshalledproducts, 5*time.Minute)
		if err != nil {
			return products, err
		}
	} else {
		err := json.Unmarshal([]byte(data), &products)
		if err != nil {
			return products, err
		}
	}
	return products, nil
}

func (r *productRepository) FindAllProductVariant(page int) ([]entity.Products, error) {
	var products []entity.Products
	key := fmt.Sprintf("FindAllProductsVariant_page_%d", page)
	const pageSize = 100

	data, _ := r.cacheable.Get(key)
	if data == "" {
		offset := (page - 1) * pageSize
		if err := r.db.Limit(pageSize).Offset(offset).Find(&products).Error; err != nil {
			return products, err
		}
		marshalledproducts, _ := json.Marshal(products)
		err := r.cacheable.Set(key, marshalledproducts, 5*time.Minute)
		if err != nil {
			return products, err
		}
	} else {
		err := json.Unmarshal([]byte(data), &products)
		if err != nil {
			return products, err
		}
	}
	return products, nil
}

func (r *productRepository) SearchByName(name string) ([]entity.Products, error) {
	var products []entity.Products
	query := r.db
	// Gunakan fungsi LOWER untuk mengabaikan perbedaan huruf besar dan kecil
	if err := query.Where("LOWER(productname) LIKE LOWER(?)", "%"+name+"%").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Update(product *entity.Products) error {
	r.cacheable.Delete("FindAllProducts_page_1")
	return r.db.Save(product).Error
}
