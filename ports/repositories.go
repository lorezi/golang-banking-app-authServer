package ports

import (
	"github.com/lorezi/golang-bank-app-auth/domain"
	"github.com/lorezi/golang-bank-app-auth/errs"
)

type AuthRepository interface {
	FindByUsernameAndPassword(username string, password string) (*domain.Login, *errs.AppError)
	StoreToken(user *domain.Login) (string, string, *errs.AppError)
	FindByToken(token string) *errs.AppError
	FindPermissionByRole(role string) ([]domain.Permission, *errs.AppError)
}
