package repository

import (
	"Wallet-System-Backend/internal/entity"
	"fmt"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTrx(trx *gorm.DB, Transaction entity.Transactions) error
	Create(transaction entity.Transactions) error
	GetAll() ([]entity.Transactions, error)
	GetByID(id uint64) (entity.Transactions, error)
	GetByUserID(id_user uint64) ([]entity.Transactions, error)
	Delete(id uint64) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTrx(trx *gorm.DB, transaction entity.Transactions) error {
	return trx.Create(&transaction).Error
}
func (r *transactionRepository) Create(transaction entity.Transactions) error {
	return r.db.Create(&transaction).Error
}

func (r *transactionRepository) GetAll() ([]entity.Transactions, error) {
	var transactions []entity.Transactions
	if err := r.db.Table("transactions").
		Select("transactions.*, users.username AS username, products.name AS product_name").
		Joins("LEFT JOIN users ON users.id = transactions.users_id").
		Joins("LEFT JOIN products ON products.id = transactions.products_id").
		Find(&transactions).Error; err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (r *transactionRepository) GetByID(id uint64) (entity.Transactions, error) {
	var Transaction entity.Transactions
	if err := r.db.First(&Transaction, id).Error; err != nil {
		return Transaction, err
	}
	return Transaction, nil
}

func (r *transactionRepository) GetByUserID(id_user uint64) ([]entity.Transactions, error) {
	var transactions []entity.Transactions
	if err := r.db.Table("transactions").
		Select("transactions.*, users.username AS username, products.name AS product_name").
		Joins("LEFT JOIN users ON users.id = transactions.users_id").
		Joins("LEFT JOIN products ON products.id = transactions.products_id").
		Where("transactions.users_id = ?", id_user).
		Find(&transactions).Error; err != nil {
		return transactions, err
	}
	return transactions, nil
}

func (r *transactionRepository) Delete(id uint64) error {
	result := r.db.Delete(&entity.Transactions{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("data with ID %d not found", id)
	}
	return result.Error
}
