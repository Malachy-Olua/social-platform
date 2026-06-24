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
