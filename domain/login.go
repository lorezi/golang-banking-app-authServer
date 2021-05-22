package domain

import (
	"database/sql"

	"github.com/dgrijalva/jwt-go"
)

type Login struct {
	Username   string         `db:"username"`
	CustomerId sql.NullString `db:"customer_id"`
	Accounts   sql.NullString `db:"account_numbers"`
	Role       string         `db:"role"`
}

type SignedDetails struct {
	Username   string
	CustomerId string
	Accounts   string
	Role       string
	jwt.StandardClaims
}
