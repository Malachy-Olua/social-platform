package store

import (
	"context"
	"database/sql"
)

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	q := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at, username, email;`

	err := s.db.QueryRowContext(ctx, q, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Username, &user.Email)

	if err != nil {
		return err
	}
	return nil

}

func (s *UsersStore) GetUserById(ctx context.Context, id int64) (*User, error) {
	q := `SELECT id, created_at, updated_at, username, email FROM users WHERE id = $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(ctx, q, id).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Username, &user.Email)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, err
		}
	}
	return &user, nil
}
