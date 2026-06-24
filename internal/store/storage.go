package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("Resource not found")

type Post struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    uuid.UUID `json:"user_id"` // ← was int64
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"` // ← had a space: json: "id"
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // ← "-" not "_" to hide from JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetPostById(ctx context.Context, id string) (*Post, error)
		// FindAll(ctx context.Context) ([]*Post, error)
	}

	Users interface {
		Create(ctx context.Context, user *User) error
		// FindById(ctx context.Context, id int) (*User, error)
		// FindAll(ctx context.Context) ([]*User, error)
	}

	Comments interface {
		GetCommentsByPostId(ctx context.Context, postId string) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db: db},
		Users:    &UsersStore{db: db},
		Comments: &CommentsStore{db: db},
	}
}
