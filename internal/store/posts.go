package store

import (
	"context"
	"database/sql"

	"time"

	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostsStore struct {
	db *sql.DB
}

type PostWithComments struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    uuid.UUID `json:"user_id"` // ← was int64
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	q := `INSERT INTO posts (title, content, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at, user_id, title, content, tags;`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, q, post.Title, post.Content, post.UserID, pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.UserID, &post.Title, &post.Content, pq.Array(&post.Tags))

	if err != nil {
		return err
	}
	return nil
}

func (s *PostsStore) GetPostById(ctx context.Context, id string) (*Post, error) {
	q := `SELECT id, created_at, updated_at, user_id, title, content, tags,version FROM posts WHERE id = $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post

	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.UserID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.Version,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *CommentsStore) GetPostWithComments(ctx context.Context) ([]*Post, error) {
	return nil, nil
}

func (s *PostsStore) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM posts WHERE id = $1;`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostsStore) UpdatePost(ctx context.Context, post *Post) error {
	query := `UPDATE posts SET title = $1, content = $2, tags = $3, version = version + 1 WHERE id = $4 AND version = $5 RETURNING version;`
	// res, err := s.db.ExecContext(ctx, query, post.Title, post.Content, pq.Array(post.Tags), post.ID, post.Version)
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, pq.Array(post.Tags), post.ID, post.Version).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound

		default:
			return err

		}
	}
	// rows, err := res.RowsAffected()
	// if err != nil {
	// 	return err
	// }
	// if rows == 0 {
	// 	return ErrNotFound
	// }

	return nil

}
