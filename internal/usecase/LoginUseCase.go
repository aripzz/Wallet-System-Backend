package usecase

import (
	"Wallet-System-Backend/internal/entity"
	"Wallet-System-Backend/internal/repository"
	"Wallet-System-Backend/utils"

	"golang.org/x/crypto/bcrypt"
)

type LoginUsecase interface {
	Authenticate(entity.UserLogin) (string, error)
}

type loginUsecase struct {
	repo repository.UserRepository
}

func NewLoginUsecase(repo repository.UserRepository) LoginUsecase {
	return &loginUsecase{repo: repo}
}

func (p *loginUsecase) Authenticate(user entity.UserLogin) (string, error) {
	foundUser, err := p.repo.FindUser(user.Username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		return "", nil
	}
	token, err := utils.GenerateJWT(user.Username, foundUser.ID, utils.GetJWTSecret())
	if err != nil {
		return "", err
	}
	return token, nil
}
