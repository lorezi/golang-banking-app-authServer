package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const ACCESS_TOKEN_DURATION = time.Hour * 24

type AccessTokenClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type AccessTokenResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
