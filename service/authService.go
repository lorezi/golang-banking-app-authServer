package service

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
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

	user, err := s.repo.FindByUsernameAndPassword(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	// verify password

	token, refreshToken, err := s.repo.StoreToken(user)
	if err != nil {
		logger.Error("error while generating and storing token")

		return nil, errs.AuthenticationError("invalid credential", "authentication failure")
	}

	return &dto.LoginResponse{AccessToken: token, RefreshToken: refreshToken}, nil
}

func (s AuthService) Verify(urlParams map[string]string) *errs.AppError {

	// RBAC
	// 1. the verify function will return the user role

	user, err := utils.Verify(urlParams["token"])
	if err != nil {
		logger.Error("error while verifying the token")
		return errs.AuthenticationError("invalid token", "authentication failure")
	}

	// 2. use the user role to query the permissions table
	permissions, appErr := s.repo.FindPermissionByRole(user.Role)
	if appErr != nil {
		logger.Error("error while verifying permissions")
		return errs.PermissionError("you don't have the right permission. Pls contact administrator", "permission failure")
	}

	// 3. loop through the returned result to check if the url params has the name in the permission and then break
	permission := false
	for _, v := range permissions {
		if strings.ToLower(urlParams["routName"]) == v.Name {
			permission = true
			break
		}

	}

	if permission {
		return nil
	}

	return errs.PermissionError("you don't have the right permission. Pls contact administrator", "permission failure")

}

func (s AuthService) Refresh(req dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError) {

	// 1. validate the token not the token expiration time
	err := utils.IsAccessTokenValid(req.AccessToken)
	if err != nil {
		return nil, errs.AuthenticationError("cannot generate a new access token until the current one expires", "access token generation failure")
	}

	// can only generate a new access token when the previous token has expired
	if err.Errors == jwt.ValidationErrorExpired {
		// 2. find token in database if it exist
		appErr := s.repo.FindByToken(req.RefreshToken)
		if appErr != nil {
			return nil, appErr
		}

		// 3. generate a new access token from refresh token
		accessToken, appErr := utils.NewAccessTokenFromRefreshToken(req.RefreshToken)
		if appErr != nil {
			return nil, appErr
		}
		return &dto.LoginResponse{AccessToken: accessToken}, nil
	}

	return nil, errs.AuthenticationError("invalid token", "new access token generation failure")

}
