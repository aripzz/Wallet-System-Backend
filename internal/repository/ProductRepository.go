package repository

import (
	"Wallet-System-Backend/internal/entity"
	"fmt"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(Product entity.CreateProducts) error
	GetAll() ([]entity.Products, error)
	GetByID(id uint64) (entity.Products, error)
	Update(Product entity.UpdateProducts) error
	UpdateTrx(trx *gorm.DB, Product entity.UpdateProducts) error
	Delete(id uint64) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(Product entity.CreateProducts) error {
	newProduct := entity.Products{
		Name:     Product.Name,
		Price:    Product.Price,
		Quantity: Product.Quantity,
	}
	return r.db.Create(&newProduct).Error
}

func (r *productRepository) GetAll() ([]entity.Products, error) {
	var Products []entity.Products
	if err := r.db.Find(&Products).Error; err != nil {
		return nil, err
	}
	return Products, nil
}

func (r *productRepository) GetByID(id uint64) (entity.Products, error) {
	var Product entity.Products
	if err := r.db.First(&Product, id).Error; err != nil {
		return Product, err
	}
	return Product, nil
}

func (r *productRepository) Update(Product entity.UpdateProducts) error {
	return r.db.Model(&Product).Where("id = ?", Product.ID).Updates(Product).Error
}

func (r *productRepository) UpdateTrx(trx *gorm.DB, Product entity.UpdateProducts) error {
	return trx.Model(&entity.Products{}).Where("id = ?", Product.ID).Updates(Product).Error
}

func (r *productRepository) Delete(id uint64) error {
	result := r.db.Delete(&entity.Products{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("data with ID %d not found", id)

	}
	return result.Error
}
