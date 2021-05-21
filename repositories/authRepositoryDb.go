// This package implements the port repositories interface
package repositories

import (
	"database/sql"

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

	// SELECT username, customer_id, role FROM users u WHERE username = ? and password = ?

	// username, u.customer_id, role

	sqlVerify := `

	SELECT username, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u
                  LEFT JOIN accounts a ON a.customer_id = u.customer_id
                WHERE username = ? and password = ?
                GROUP BY a.customer_id	
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

/* generate and store the token */
func (r AuthRepositoryDb) StoreToken() (string, *errs.AppError) {
	// generate token
	token, appErr := utils.GenerateJwt()
	if appErr != nil {
		return "", appErr
	}

	// save it in the store
	sqlInsert := "INSERT into refresh_token_store (refresh_token) values (?)"

	_, err := r.client.Exec(sqlInsert, token)
	if err != nil {
		logger.Error("unexpected database error: " + err.Error())

		return "", errs.UnExpectedError("unexpected database error", "error")
	}

	return token, nil
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
