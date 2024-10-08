package repository

import (
	"Wallet-System-Backend/internal/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUser(username string) (entity.Users, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindUser(username string) (entity.Users, error) {
	var user entity.Users
	if err := r.db.Where(&entity.Users{Username: username}).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
