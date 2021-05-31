// This package implements the port repositories interface
package repositories

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lorezi/golang-bank-app-auth/domain"
	"github.com/lorezi/golang-bank-app-auth/errs"
	"github.com/lorezi/golang-bank-app-auth/logger"
	"github.com/lorezi/golang-bank-app-auth/utils"
)

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func NewAuthRepositoryDb(client *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}

// find user by username and password
func (r AuthRepositoryDb) FindByUsernameAndPassword(username string, password string) (*domain.Login, *errs.AppError) {

	login := &domain.Login{}

	sqlVerify := `

	SELECT username, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u
                  LEFT JOIN accounts a ON a.customer_id = u.customer_id
                WHERE username = ? and password = ?
                GROUP BY username, u.customer_id, role	
	`

	err := r.client.Get(login, sqlVerify, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.AuthenticationError("invalid credentials", "fail")
		}
		logger.Error("error while verifying login request from the database: " + err.Error())
		return nil, errs.UnExpectedError("unexpected database error", "error")
	}

	return login, nil

}

// generate and store the token - returns token, refreshToken and error
func (r AuthRepositoryDb) StoreToken(user *domain.Login) (string, string, *errs.AppError) {
	// generate token
	tokenClaims := user.GenerateTokenClaims()
	token, appErr := utils.GenerateJwt(&tokenClaims)
	if appErr != nil {
		return "", "", appErr
	}

	// the refresh access token is use to generate a new access token once
	refreshTokenClaims := user.GenerateRefreshTokenClaims()
	refreshToken, appErr := utils.GenerateRefreshJwt(&refreshTokenClaims)
	if appErr != nil {
		return "", "", appErr
	}

	// save it in the store
	sqlInsert := "INSERT into refresh_token_store (refresh_token) values (?)"

	_, err := r.client.Exec(sqlInsert, refreshToken)
	if err != nil {
		logger.Error("unexpected database error: " + err.Error())
		return "", "", errs.UnExpectedError("unexpected database error", "error")
	}

	return token, refreshToken, nil
}

// refresh the token if it exists
func (r AuthRepositoryDb) FindByToken(token string) *errs.AppError {

	sqlSelect := "SELECT refresh_token FROM refresh_token_store WHERE refresh_token = ?"

	dbToken := ""
	err := r.client.Get(&dbToken, sqlSelect, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.AuthenticationError("refresh token not registered in the store", "fail")
		} else {
			logger.Error("Unexpected database error: " + err.Error())
			return errs.UnExpectedError("unexpected database error", "error")
		}
	}

	return nil
}

func (r AuthRepositoryDb) FindPermissionByRole(role string) ([]domain.Permission, *errs.AppError) {

	sp := []domain.Permission{}

	sqlSelect := fmt.Sprintf("SELECT permission_name FROM permissions WHERE role_name = '%s'", role)

	err := r.client.Select(&sp, sqlSelect)

	if err != nil {
		logger.Error("Error while scanning permissions " + err.Error())
		return nil, errs.UnExpectedError("unexpected database error", "error")
	}

	return sp, nil

}
