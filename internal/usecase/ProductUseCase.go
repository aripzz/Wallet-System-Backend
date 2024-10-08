package usecase

import (
	"Wallet-System-Backend/infra"
	"Wallet-System-Backend/internal/constant"
	"Wallet-System-Backend/internal/entity"
	"Wallet-System-Backend/internal/repository"
	"encoding/json"
	"errors"
	"strconv"
)

type ProductsUsecase interface {
	Create(product entity.CreateProducts) error
	GetAll() ([]entity.Products, error)
	GetByID(id uint64) (entity.Products, error)
	Update(product entity.UpdateProducts) error
	Delete(id uint64) error
}

type productUsecase struct {
	repo  repository.ProductRepository
	cache *infra.RedisClient
}

func NewProductsUsecase(repo repository.ProductRepository, cache *infra.RedisClient) ProductsUsecase {
	return &productUsecase{repo: repo, cache: cache}
}

func (p *productUsecase) Create(product entity.CreateProducts) error {
	if product.Price == 0 {
		return errors.New(constant.ErrBadRequest + "price cannot be zero")
	}
	if product.Quantity <= 0 {
		return errors.New(constant.ErrBadRequest + "quantity cannot be zero or negative")
	}
	err := p.repo.Create(product)
	if err != nil {
		return err
	}

	p.cache.Delete("products")
	return nil
}

func (p *productUsecase) GetAll() ([]entity.Products, error) {
	cachedProductss, err := p.cache.Get("products")
	if err == nil && cachedProductss != "" {
		var products []entity.Products
		json.Unmarshal([]byte(cachedProductss), &products)
		return products, nil
	}

	products, err := p.repo.GetAll()
	if err != nil {
		return nil, err
	}

	cachedData, _ := json.Marshal(products)
	p.cache.Set("products", string(cachedData))
	return products, nil
}

func (p *productUsecase) GetByID(id uint64) (entity.Products, error) {
	cachedProducts, err := p.cache.Get("product:" + strconv.Itoa(int(id)))
	if err == nil && cachedProducts != "" {
		var product entity.Products
		json.Unmarshal([]byte(cachedProducts), &product)
		return product, nil
	}

	product, err := p.repo.GetByID(id)
	if err != nil {
		return product, err
	}

	cachedData, _ := json.Marshal(product)
	p.cache.Set("product:"+strconv.Itoa(int(id)), string(cachedData))
	return product, nil
}

func (p *productUsecase) Update(product entity.UpdateProducts) error {
	existingProducts, err := p.repo.GetByID(product.ID)
	if err != nil {
		return err
	}

	if existingProducts.ID == 0 {
		return errors.New("not found")
	}

	err = p.repo.Update(product)
	if err != nil {
		return err
	}

	p.cache.Delete("product:" + strconv.Itoa(int(product.ID)))
	p.cache.Delete("products")

	return nil
}

func (p *productUsecase) Delete(id uint64) error {
	err := p.repo.Delete(id)
	if err != nil {
		return err
	}

	p.cache.Delete("product:" + strconv.Itoa(int(id)))
	p.cache.Delete("products")
	return nil
}
