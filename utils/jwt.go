package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lorezi/golang-bank-app-auth/errs"
	"github.com/lorezi/golang-bank-app-auth/logger"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateJwt() (string, *errs.AppError) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    "user",
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	token, err := claims.SignedString([]byte(SECRET_KEY))

	if err != nil {
		logger.Error("Failed while signing refresh token: " + err.Error())
		return "", errs.UnExpectedError("cannot generate refresh token", "fail")
	}

	return token, nil
}

func Verify(token string) error {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil || !jwtToken.Valid {
		return err
	}

	// claims := jwtToken.Claims.(*jwt.StandardClaims)
	_ = jwtToken.Claims.(*jwt.StandardClaims)

	return nil
}
