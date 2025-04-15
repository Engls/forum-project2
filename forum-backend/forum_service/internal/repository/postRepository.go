package repository

import (
	"context"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error)
	GetPosts(ctx context.Context) ([]entity.Post, error)
	GetPostByID(ctx context.Context, id int) (*entity.Post, error)
	UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error)
	DeletePost(ctx context.Context, id int) error
	GetUserIDByToken(ctx context.Context, token string) (int, error)
}

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	query := `INSERT INTO posts (author_id, title, content) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, post.AuthorId, post.Title, post.Content)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	post.ID = int(id)
	return &post, nil
}

func (r *postRepository) GetPosts(ctx context.Context) ([]entity.Post, error) {
	query := `SELECT id, author_id, title, content FROM posts`
	var posts []entity.Post
	err := r.db.SelectContext(ctx, &posts, query)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepository) GetPostByID(ctx context.Context, id int) (*entity.Post, error) {
	query := `SELECT id, author_id, title, content FROM posts WHERE id = ?`
	var post entity.Post
	err := r.db.GetContext(ctx, &post, query, id)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	query := `UPDATE posts SET title = ?, content = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) DeletePost(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *postRepository) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	query := `SELECT user_id FROM tokens WHERE token = ?`
	var userID int
	err := r.db.GetContext(ctx, &userID, query, token)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
