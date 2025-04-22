package repository

import (
	"context"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"go.uber.org/zap"
)

type DBposts interface {
	Get(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type PostRepository interface {
	CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error)
	GetPosts(ctx context.Context) ([]entity.Post, error)
	GetPostByID(ctx context.Context, id int) (*entity.Post, error)
	UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error)
	DeletePost(ctx context.Context, id int) error
	GetUserIDByToken(ctx context.Context, token string) (int, error)
}

type postRepository struct {
	db     DB
	logger *zap.Logger
}

func NewPostRepository(db DB, logger *zap.Logger) PostRepository {
	return &postRepository{db: db, logger: logger}
}

func (r *postRepository) CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	query := `INSERT INTO posts (author_id, title, content) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, post.AuthorId, post.Title, post.Content)
	if err != nil {
		r.logger.Error("Failed to create post", zap.Error(err), zap.Int("authorID", post.AuthorId))
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.logger.Error("Failed to get last insert ID", zap.Error(err))
		return nil, err
	}

	post.ID = int(id)
	r.logger.Info("Post created successfully", zap.Int("postID", post.ID), zap.Int("authorID", post.AuthorId))
	return &post, nil
}

func (r *postRepository) GetPosts(ctx context.Context) ([]entity.Post, error) {
	query := `SELECT id, author_id, title, content FROM posts`
	var posts []entity.Post
	err := r.db.SelectContext(ctx, &posts, query)
	if err != nil {
		r.logger.Error("Failed to get posts", zap.Error(err))
		return nil, err
	}
	r.logger.Info("Posts retrieved successfully", zap.Int("count", len(posts)))
	return posts, nil
}

func (r *postRepository) GetPostByID(ctx context.Context, id int) (*entity.Post, error) {
	query := `SELECT id, author_id, title, content FROM posts WHERE id = ?`
	var post entity.Post
	err := r.db.GetContext(ctx, &post, query, id)
	if err != nil {
		r.logger.Error("Failed to get post by ID", zap.Error(err), zap.Int("postID", id))
		return nil, err
	}
	r.logger.Info("Post retrieved successfully", zap.Int("postID", id))
	return &post, nil
}

func (r *postRepository) UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	query := `UPDATE posts SET title = ?, content = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		r.logger.Error("Failed to update post", zap.Error(err), zap.Int("postID", post.ID))
		return nil, err
	}
	r.logger.Info("Post updated successfully", zap.Int("postID", post.ID))
	return &post, nil
}

func (r *postRepository) DeletePost(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete post", zap.Error(err), zap.Int("postID", id))
		return err
	}
	r.logger.Info("Post deleted successfully", zap.Int("postID", id))
	return nil
}

func (r *postRepository) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	query := `SELECT user_id FROM tokens WHERE token = ?`
	var userID int
	err := r.db.GetContext(ctx, &userID, query, token)
	if err != nil {
		r.logger.Error("Failed to get user ID by token", zap.Error(err), zap.String("token", token))
		return 0, err
	}
	r.logger.Info("User ID retrieved successfully", zap.String("token", token), zap.Int("userID", userID))
	return userID, nil
}
