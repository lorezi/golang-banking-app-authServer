package utils

import (
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/lorezi/golang-bank-app-auth/domain"
	"github.com/lorezi/golang-bank-app-auth/errs"
	"github.com/lorezi/golang-bank-app-auth/logger"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateJwt(tokenClaims *domain.AccessTokenClaims) (string, *errs.AppError) {
	return newToken(tokenClaims)
}

func GenerateRefreshJwt(refreshTokenClaims *domain.RefreshTokenClaims) (string, *errs.AppError) {
	return newToken(refreshTokenClaims)
}

func Verify(token string) (*domain.AccessTokenResponse, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &domain.AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil || !jwtToken.Valid {
		return nil, err
	}

	claims := jwtToken.Claims.(*domain.AccessTokenClaims)
	claimsDetail := &domain.AccessTokenResponse{
		Username: claims.Username,
		Role:     claims.Role,
	}

	return claimsDetail, nil
}

func IsAccessTokenValid(token string) *jwt.ValidationError {
	// 1. checks the validity of the token not the expiration time
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		vErr := &jwt.ValidationError{}
		if errors.As(err, &vErr) {
			return vErr
		}
	}

	return nil
}

func NewAccessTokenFromRefreshToken(refreshToken string) (string, *errs.AppError) {
	token, err := jwt.ParseWithClaims(refreshToken, &domain.AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil || !token.Valid {
		return "", errs.UnExpectedError(err.Error(), "fail")
	}

	tokenClaims := token.Claims.(*domain.AccessTokenClaims)

	return newToken(tokenClaims)

}

func newToken(tokenClaims jwt.Claims) (string, *errs.AppError) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	token, err := claims.SignedString([]byte(SECRET_KEY))

	if err != nil {
		logger.Error("Failed while signing refresh token: " + err.Error())
		return "", errs.UnExpectedError("cannot generate refresh token", "fail")
	}

	return token, nil
}
