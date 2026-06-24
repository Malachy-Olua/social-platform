package store

import (
	"context"
	"database/sql"

	"time"

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

	err := s.db.QueryRowContext(ctx, q, post.Title, post.Content, post.UserID, pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.UserID, &post.Title, &post.Content, pq.Array(&post.Tags))

	if err != nil {
		return err
	}
	return nil
}

func (s *PostsStore) GetPostById(ctx context.Context, id string) (*Post, error) {
	q := `SELECT id, created_at, updated_at, user_id, title, content, tags FROM posts WHERE id = $1;`
	var post Post

	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.UserID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
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
