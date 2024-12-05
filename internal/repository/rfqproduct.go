package repository

import (
	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"gorm.io/gorm"
)

type RfqProductRepository interface {
	CreateProduct(product *entity.RfqsProduct) (*entity.RfqsProduct, error)
	GetLastProductId() (string, error)
	GetProductsByRfqId(rfqId string) ([]entity.RfqsProduct, error)
	CheckMaterialExists(materialId string) (bool, error)
	UpdateProduct(product *entity.RfqsProduct) (*entity.RfqsProduct, error)
	GetProductByRfqIdAndProductId(rfqId, productId string) (*entity.RfqsProduct, error)
	GetProductDetails(materialId string) (*entity.Materials, error)
}

type rfqProductRepository struct {
	db *gorm.DB
}

func NewRfqProductRepository(db *gorm.DB) RfqProductRepository {
	return &rfqProductRepository{db: db}
}

func (r *rfqProductRepository) CreateProduct(product *entity.RfqsProduct) (*entity.RfqsProduct, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *rfqProductRepository) GetLastProductId() (string, error) {
	var lastProduct entity.RfqsProduct
	err := r.db.Order("id_rfqproduct DESC").First(&lastProduct).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", err
	}
	return lastProduct.RfqsProductId, nil
}

func (r *rfqProductRepository) GetProductsByRfqId(rfqId string) ([]entity.RfqsProduct, error) {
	var Products []entity.RfqsProduct
	if err := r.db.Where("id_rfq = ?", rfqId).Find(&Products).Error; err != nil {
		return nil, err
	}
	return Products, nil
}

func (r *rfqProductRepository) CheckMaterialExists(materialId string) (bool, error) {
	var count int64
	if err := r.db.Model(&entity.Materials{}).Where("id_material = ?", materialId).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *rfqProductRepository) UpdateProduct(product *entity.RfqsProduct) (*entity.RfqsProduct, error) {
	// Perbarui produk berdasarkan id_rfqproduct atau kombinasi id_rfq dan id_product
	if err := r.db.Model(&entity.RfqsProduct{}).
		Where("id_rfqproduct = ? AND deleted_at IS NULL", product.RfqsProductId).
		Updates(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *rfqProductRepository) GetProductByRfqIdAndProductId(rfqId, productId string) (*entity.RfqsProduct, error) {
	var product entity.RfqsProduct
	if err := r.db.Where("id_rfq = ? AND id_material = ?", rfqId, productId).First(&product).Error; err != nil {
		// Jika tidak ditemukan, kembalikan nil dan error yang sesuai
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *rfqProductRepository) GetProductDetails(materialId string) (*entity.Materials, error) {
	var product entity.Materials
	if err := r.db.Table("material").Where("id_material = ?", materialId).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
