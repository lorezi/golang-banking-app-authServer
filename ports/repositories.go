package ports

import (
	"github.com/lorezi/golang-bank-app-auth/domain"
	"github.com/lorezi/golang-bank-app-auth/errs"
)

type AuthRepository interface {
	FindByUsernameAndPassword(username string, password string) (*domain.Login, *errs.AppError)
	StoreToken() (string, *errs.AppError)
	FindByToken(token string) *errs.AppError
}
