package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type CommentsStore struct {
	db *sql.DB
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

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	return nil
}

func (s *CommentsStore) GetCommentsByPostId(ctx context.Context, postId string) ([]Comment, error) {
	q :=
		`SELECT 
		comments.id, 
		post_id, 
		comments.user_id, 
		users.username,
		users.id,
		content, 
		comments.created_at, 
		comments.updated_at 
		FROM comments 
		JOIN users ON users.id = comments.user_id 
		WHERE post_id = $1
		ORDER BY comments.created_at DESC;
	`

	rows, err := s.db.QueryContext(ctx, q, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		comment.User = User{}
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.User.Username, &comment.User.ID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
