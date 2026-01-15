package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserStore interface {
	GetByEmail(email string) (*User, error)
}

type User struct {
	Id          string    `json:"id"`
	Email       string    `json:"email"`
	IsActivated bool      `json:"is_activated"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, is_activated, created_at, updated_at
		FROM users
		WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.IsActivated,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
