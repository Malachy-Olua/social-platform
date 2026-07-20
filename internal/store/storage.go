package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("Resource not found")
var QueryTimeoutDuration = time.Second * 5

type Post struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    uuid.UUID `json:"user_id"` // ← was int64
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
}

type User struct {
	ID        uuid.UUID `json:"id"` // ← had a space: json: "id"
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // ← "-" not "_" to hide from JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	ID        uuid.UUID `json:"id"`
	PostID    uuid.UUID `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User
}

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetPostById(ctx context.Context, id string) (*Post, error)
		Delete(ctx context.Context, id string) error
		UpdatePost(ctx context.Context, post *Post) error

		// FindAll(ctx context.Context) ([]*Post, error)
	}

	Users interface {
		Create(ctx context.Context, user *User) error
		GetUserById(ctx context.Context, id int64) (*User, error)
		// FindAll(ctx context.Context) ([]*User, error)
	}

	Comments interface {
		Create(ctx context.Context, Comment *Comment) error
		GetCommentsByPostId(ctx context.Context, postId string) ([]Comment, error)
	}

	Followers interface {
		FollowUser(ctx context.Context, followerId string, followeeId string) error
		UnfollowUser(ctx context.Context, followerId string, followeeId string) error
		GetFollowers(ctx context.Context, userId string) ([]User, error)
		GetFollowing(ctx context.Context, userId string) ([]User, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostsStore{db: db},
		Users:     &UsersStore{db: db},
		Comments:  &CommentsStore{db: db},
		Followers: &FollowersStore{db: db},
	}
}
