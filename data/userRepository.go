package data

import "database/sql"

type UserRepository struct {
	DB *sql.DB
}
