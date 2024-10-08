package usecase

import (
	"Wallet-System-Backend/infra"
	"Wallet-System-Backend/internal/constant"
	"Wallet-System-Backend/internal/entity"
	"Wallet-System-Backend/internal/repository"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type TransactionUseCase interface {
	CreateTransaction(userID uint64, transaction entity.RequestCreateTransactions) error
	GetAllTransactions() ([]entity.Transactions, error)
	GetTransactionByID(id uint64) (entity.Transactions, error)
	GetTransactionsByUserID(id_user uint64) ([]entity.Transactions, error)
	DeleteTransaction(id uint64) error
}

type transactionUseCase struct {
	transactionRepo repository.TransactionRepository
	productRepo     repository.ProductRepository
	walletRepo      repository.WalletRepository
	db              *gorm.DB
	cache           *infra.RedisClient
}

func NewTransactionUseCase(transactionRepo repository.TransactionRepository, productRepo repository.ProductRepository, walletRepo repository.WalletRepository, db *gorm.DB, cache *infra.RedisClient) TransactionUseCase {
	return &transactionUseCase{transactionRepo: transactionRepo, productRepo: productRepo, walletRepo: walletRepo, db: db, cache: cache}
}

func (u *transactionUseCase) CreateTransaction(userID uint64, transaction entity.RequestCreateTransactions) error {
	trx := u.db.Begin()
	var err error
	var amount float64
	defer func() {
		if err != nil {
			trx.Rollback()
			u.transactionRepo.Create(entity.Transactions{
				UsersID:    userID,
				ProductsID: transaction.ProductID,
				Amount:     amount,
				Status:     "failed",
			})
		} else {
			trx.Commit()
		}
	}()

	product, err := u.productRepo.GetByID(transaction.ProductID)
	if err != nil {
		return err
	}

	if product.Quantity < 1 {
		return gorm.ErrRecordNotFound
	}

	product.Quantity--

	err = u.productRepo.UpdateTrx(trx, entity.UpdateProducts{
		ID:       product.ID,
		Quantity: &product.Quantity,
	})
	if err != nil {
		return err
	}

	amount = product.Price

	wallet, err := u.walletRepo.GetByWalletTypes(transaction.WalletTypeID, userID)
	if err != nil {
		return err
	}

	if wallet.Balance < amount {
		err = errors.New(constant.ErrBadRequest + "Saldo tidak cukup untuk melakukan pembelian. Saldo saat ini: " + fmt.Sprintf("%d", int64(wallet.Balance)))
		return err
	}

	if !wallet.Active {
		err = errors.New(constant.ErrBadRequest + "wallet is not active")
		return err
	}

	newBalence := wallet.Balance - amount

	err = u.walletRepo.UpdateTrx(trx, entity.UpdateWallets{
		ID:      wallet.ID,
		Balance: &newBalence,
	})
	if err != nil {
		return err
	}

	err = u.transactionRepo.CreateTrx(trx, entity.Transactions{
		UsersID:    userID,
		ProductsID: transaction.ProductID,
		Amount:     amount,
		Status:     "success",
	})

	u.cache.Delete("wallet:" + strconv.Itoa(int(wallet.ID)))
	u.cache.Delete("wallets")

	return err
}

func (u *transactionUseCase) GetAllTransactions() ([]entity.Transactions, error) {
	return u.transactionRepo.GetAll()
}

func (u *transactionUseCase) GetTransactionByID(id uint64) (entity.Transactions, error) {
	return u.transactionRepo.GetByID(id)
}

func (u *transactionUseCase) GetTransactionsByUserID(id_user uint64) ([]entity.Transactions, error) {
	return u.transactionRepo.GetByUserID(id_user)
}

func (u *transactionUseCase) DeleteTransaction(id uint64) error {
	return u.transactionRepo.Delete(id)
}
