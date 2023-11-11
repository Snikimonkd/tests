package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"lab1/internal/model"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) repository {
	return repository{
		db: db,
	}
}

func (r repository) CheckUserExist(email string) (bool, error) {
	query := `SELECT email FROM users WHERE email = $1;`

	var buf *string
	err := r.db.QueryRow(query, email).Scan(&buf)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("can't check if user alreay exist: %w", err)
	}

	return true, nil
}

func (r repository) CreateUser(user model.User) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2);`

	_, err := r.db.Exec(query, user.Name, user.Email)
	if err != nil {
		return fmt.Errorf("can't insert into users: %w", err)
	}

	return nil
}
