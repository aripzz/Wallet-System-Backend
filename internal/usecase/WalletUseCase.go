package usecase

import (
	"Wallet-System-Backend/infra"
	"Wallet-System-Backend/internal/entity"
	"Wallet-System-Backend/internal/repository"
	"encoding/json"
	"errors"
	"strconv"
)

type WalletsUsecase interface {
	Create(wallet entity.CreateWallets) error
	GetAll() ([]entity.Wallets, error)
	GetByID(id uint64) (entity.Wallets, error)
	Update(wallet entity.UpdateWallets) error
	Delete(id uint64) error
}

type walletUsecase struct {
	repo  repository.WalletRepository
	cache *infra.RedisClient
}

func NewWalletsUsecase(repo repository.WalletRepository, cache *infra.RedisClient) WalletsUsecase {
	return &walletUsecase{repo: repo, cache: cache}
}

func (p *walletUsecase) Create(wallet entity.CreateWallets) error {
	err := p.repo.Create(wallet)
	if err != nil {
		return err
	}

	p.cache.Delete("wallets")
	return nil
}

func (p *walletUsecase) GetAll() ([]entity.Wallets, error) {
	cachedWalletss, err := p.cache.Get("wallets")
	if err == nil && cachedWalletss != "" {
		var wallets []entity.Wallets
		json.Unmarshal([]byte(cachedWalletss), &wallets)
		return wallets, nil
	}

	wallets, err := p.repo.GetAll()
	if err != nil {
		return nil, err
	}

	cachedData, _ := json.Marshal(wallets)
	p.cache.Set("wallets", string(cachedData))
	return wallets, nil
}

func (p *walletUsecase) GetByID(id uint64) (entity.Wallets, error) {
	cachedWallets, err := p.cache.Get("wallet:" + strconv.Itoa(int(id)))
	if err == nil && cachedWallets != "" {
		var wallet entity.Wallets
		json.Unmarshal([]byte(cachedWallets), &wallet)
		return wallet, nil
	}

	wallet, err := p.repo.GetByID(id)
	if err != nil {
		return wallet, err
	}

	cachedData, _ := json.Marshal(wallet)
	p.cache.Set("wallet:"+strconv.Itoa(int(id)), string(cachedData))
	return wallet, nil
}

func (p *walletUsecase) Update(wallet entity.UpdateWallets) error {
	existingWallets, err := p.repo.GetByID(wallet.ID)
	if err != nil {
		return err
	}

	if existingWallets.ID == 0 {
		return errors.New("not found")
	}

	err = p.repo.Update(wallet)
	if err != nil {
		return err
	}

	p.cache.Delete("wallet:" + strconv.Itoa(int(wallet.ID)))
	p.cache.Delete("wallets")

	return nil
}

func (p *walletUsecase) Delete(id uint64) error {
	err := p.repo.Delete(id)
	if err != nil {
		return err
	}

	p.cache.Delete("wallet:" + strconv.Itoa(int(id)))
	p.cache.Delete("wallets")
	return nil
}
