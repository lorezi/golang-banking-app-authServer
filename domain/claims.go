package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const ACCESS_TOKEN_DURATION = time.Hour * 24
const REFRESH_TOKEN_DURATION = time.Hour * 168

type AccessTokenClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	TokenType string `json:"token_type"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	jwt.StandardClaims
}

type AccessTokenResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
