package service

import (
	"github.com/lorezi/golang-bank-app-auth/dto"
	"github.com/lorezi/golang-bank-app-auth/errs"
	"github.com/lorezi/golang-bank-app-auth/logger"
	"github.com/lorezi/golang-bank-app-auth/ports"
	"github.com/lorezi/golang-bank-app-auth/utils"
)

type AuthService struct {
	repo ports.AuthRepository
}

func NewAuthService(repo ports.AuthRepository) *AuthService {
	return &AuthService{repo}
}

func (s AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {

	// login := &domain.Login{}

	_, err := s.repo.FindByUsernameAndPassword(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	// verify password

	token, err := s.repo.StoreToken()
	if err != nil {
		logger.Error("error while generating and storing token")

		return nil, errs.AuthenticationError("invalid credential", "fail")
	}

	return &dto.LoginResponse{AccessToken: token}, nil
}

func (s AuthService) Verify(urlParams map[string]string) *errs.AppError {
	err := utils.Verify(urlParams["token"])
	if err != nil {
		return errs.AuthenticationError("invalid token", "authentication failure")
	}

	return nil
}
