package repository

import (
	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"gorm.io/gorm"
)

type QuoProductRepository interface {
	GetLastProductId() (string, error)
	CreateProduct(product *entity.QuotationsProduct) (*entity.QuotationsProduct, error)
	GetProductsByQuoId(rfqId string) ([]entity.QuotationsProduct, error)
	GetProductByQuoIdAndProductId(rfqId, productId string) (*entity.QuotationsProduct, error)
	CheckMaterialExists(materialId string) (bool, error)
	UpdateProduct(product *entity.QuotationsProduct) (*entity.QuotationsProduct, error)
	GetProductDetails(materialId string) (*entity.Products, error)
	DeleteProductsByQuoId(rfqId string) error
}

type quoProductRepository struct {
	db *gorm.DB
}

func NewQuoProductRepository(db *gorm.DB) QuoProductRepository {
	return &quoProductRepository{db: db}
}

func (r *quoProductRepository) GetLastProductId() (string, error) {
	var lastProduct entity.QuotationsProduct
	err := r.db.Order("id_quotationsproduct DESC").First(&lastProduct).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	return lastProduct.QuotationsProductId, nil
}
func (r *quoProductRepository) CreateProduct(product *entity.QuotationsProduct) (*entity.QuotationsProduct, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}
func (r *quoProductRepository) GetProductsByQuoId(rfqId string) ([]entity.QuotationsProduct, error) {
	var Products []entity.QuotationsProduct
	if err := r.db.Where("id_quotation = ?", rfqId).Find(&Products).Error; err != nil {
		return nil, err
	}
	return Products, nil
}
func (r *quoProductRepository) CheckMaterialExists(materialId string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Products{}).Where("id_product = ?", materialId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *quoProductRepository) GetProductByQuoIdAndProductId(rfqId, productId string) (*entity.QuotationsProduct, error) {
	var product entity.QuotationsProduct
	if err := r.db.Where("id_quotation = ? AND id_product = ?", rfqId, productId).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *quoProductRepository) UpdateProduct(product *entity.QuotationsProduct) (*entity.QuotationsProduct, error) {
	// Perbarui produk berdasarkan id_rfqproduct atau kombinasi id_rfq dan id_product
	if err := r.db.Model(&entity.QuotationsProduct{}).
		Where("id_quotationsproduct = ? AND deleted_at IS NULL", product.QuotationsProductId).
		Updates(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *quoProductRepository) GetProductDetails(materialId string) (*entity.Products, error) {
	var product entity.Products
	if err := r.db.Table("products").Where("id_product = ?", materialId).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *quoProductRepository) DeleteProductsByQuoId(rfqId string) error {
	result := s.db.Unscoped().Where("id_quotation = ?", rfqId).Delete(&entity.QuotationsProduct{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
