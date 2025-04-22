package repository

import (
	"context"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"go.uber.org/zap"
)

type CommentsRepository interface {
	CreateComment(ctx context.Context, comment entity.Comment) (entity.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID int) ([]entity.Comment, error)
}

type commentsRepository struct {
	db     DB
	logger *zap.Logger
}

func NewCommentsRepository(db DB, logger *zap.Logger) CommentsRepository {
	return &commentsRepository{db: db, logger: logger}
}

func (r *commentsRepository) CreateComment(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	query := `
		INSERT INTO comments (post_id, author_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query, comment.PostId, comment.AuthorId, comment.Content).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		r.logger.Error("Failed to create comment", zap.Error(err), zap.Int("postID", comment.PostId), zap.Int("authorID", comment.AuthorId))
		return entity.Comment{}, err
	}
	r.logger.Info("Comment created successfully", zap.Int("commentID", comment.ID), zap.Int("postID", comment.PostId), zap.Int("authorID", comment.AuthorId))
	return comment, nil
}

func (r *commentsRepository) GetCommentsByPostID(ctx context.Context, postID int) ([]entity.Comment, error) {
	query := `
		SELECT id, post_id, author_id, content, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		r.logger.Error("Failed to get comments by post ID", zap.Error(err), zap.Int("postID", postID))
		return nil, err
	}
	defer rows.Close()

	var comments []entity.Comment
	for rows.Next() {
		var comment entity.Comment
		err := rows.Scan(&comment.ID, &comment.PostId, &comment.AuthorId, &comment.Content, &comment.CreatedAt)
		if err != nil {
			r.logger.Error("Failed to scan comment", zap.Error(err), zap.Int("postID", postID))
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		r.logger.Error("Error occurred while iterating over comments", zap.Error(err), zap.Int("postID", postID))
		return nil, err
	}
	r.logger.Info("Comments retrieved successfully", zap.Int("postID", postID), zap.Int("count", len(comments)))
	return comments, nil
}
