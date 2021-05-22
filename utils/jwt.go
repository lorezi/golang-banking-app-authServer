package utils

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/lorezi/golang-bank-app-auth/domain"
	"github.com/lorezi/golang-bank-app-auth/errs"
	"github.com/lorezi/golang-bank-app-auth/logger"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateJwt(tokenClaims *domain.AccessTokenClaims) (string, *errs.AppError) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	token, err := claims.SignedString([]byte(SECRET_KEY))

	if err != nil {
		logger.Error("Failed while signing refresh token: " + err.Error())
		return "", errs.UnExpectedError("cannot generate refresh token", "fail")
	}

	return token, nil
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

	// _ = jwtToken.Claims.(*jwt.StandardClaims)

	return claimsDetail, nil
}
