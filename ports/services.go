package ports

import (
	"github.com/lorezi/golang-bank-app-auth/dto"
	"github.com/lorezi/golang-bank-app-auth/errs"
)

type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
}
