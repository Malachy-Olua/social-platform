package store

import (
	"context"
	"database/sql"
)

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) Create(ctx context.Context, Comment *Comment) error {
	q := `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id, user_id, content, post_id, created_at, updated_at;`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, q, Comment.PostID, Comment.UserID, Comment.Content).
		Scan(&Comment.ID, &Comment.UserID, &Comment.Content, &Comment.PostID, &Comment.CreatedAt, &Comment.UpdatedAt)

	if err != nil {
		return err
	}
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
