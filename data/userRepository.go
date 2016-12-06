package data

import (
	"database/sql"
	"log"
	"fmt"
	"github.com/markbates/goth"
	"errors"
)

type UserRepository struct {
	DB *sql.DB
}

// GetUserId checks for an existing oauth user in the authorizations table.
// If none exist, it creates entries in the users table and the authorizations table.
// It returns the id from the users table.
func (r *UserRepository) GetUserId(provider string, user goth.User) (userId int64, err error) {
	err = r.DB.QueryRow(
		"SELECT user_id FROM authorizations WHERE provider LIKE ? AND provider_id = ?",
		provider,
		user.UserID,
	).Scan(&userId)
	if err != sql.ErrNoRows {
		return
	}

	// create user
	result, err := r.DB.Exec("INSERT INTO users () VALUES(?)")
	if err != nil {
		return
	}
	userId, err = result.LastInsertId()
	if err != nil {
		return
	}
	result, err = r.DB.Exec(
		"INSERT INTO authorizations (user_id, provider, provider_user_id, name, email) VALUES(?, ?, ?, ?, ?)",
		userId,
		provider,
		user.UserID,
		user.Name,
		user.Email,
	)
	if err != nil {
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rows != 1 {
		err = errors.New(fmt.Sprintf("Authorization INSERT affected %d rows", rows))
	}
	return
}