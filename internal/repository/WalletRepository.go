package repository

import (
	"Wallet-System-Backend/internal/entity"
	"fmt"

	"gorm.io/gorm"
)

type WalletRepository interface {
	Create(wallet entity.CreateWallets) error
	GetAll() ([]entity.Wallets, error)
	GetByID(id uint64) (entity.Wallets, error)
	GetByWalletTypes(wallet_types_id uint64, user_id uint64) (entity.Wallets, error)
	Update(Wallet entity.UpdateWallets) error
	Delete(id uint64) error
	UpdateTrx(trx *gorm.DB, wallet entity.UpdateWallets) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(wallet entity.CreateWallets) error {
	return r.db.Create(&wallet).Error
}

func (r *walletRepository) GetAll() ([]entity.Wallets, error) {
	var Wallets []entity.Wallets
	if err := r.db.Find(&Wallets).Error; err != nil {
		return nil, err
	}

	return Wallets, nil
}
func (r *walletRepository) GetByWalletTypes(wallet_types_id uint64, user_id uint64) (entity.Wallets, error) {
	var Wallet entity.Wallets
	if err := r.db.Where("wallet_types_id = ?", wallet_types_id).Where("user_id = ?", user_id).First(&Wallet).Error; err != nil {
		return Wallet, err
	}
	return Wallet, nil
}

func (r *walletRepository) GetByID(id uint64) (entity.Wallets, error) {
	var Wallet entity.Wallets
	if err := r.db.First(&Wallet, id).Error; err != nil {
		return Wallet, err
	}
	return Wallet, nil
}

func (r *walletRepository) Update(Wallet entity.UpdateWallets) error {
	return r.db.Model(&entity.Wallets{}).Where("id = ?", Wallet.ID).Updates(Wallet).Error
}

func (r *walletRepository) UpdateTrx(trx *gorm.DB, wallet entity.UpdateWallets) error {
	return trx.Model(&entity.Wallets{}).Where("id = ?", wallet.ID).Updates(wallet).Error

}

func (r *walletRepository) Delete(id uint64) error {
	result := r.db.Delete(&entity.Wallets{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("data with ID %d not found", id)
	}
	return result.Error
}
