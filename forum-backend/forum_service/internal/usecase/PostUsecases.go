package usecase

import (
	"context"
	"forum-project/forum_service/internal/entity"
	"forum-project/forum_service/internal/repository"
)

type PostUsecase interface {
	CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error)
	GetPosts(ctx context.Context) ([]entity.Post, error)
	GetPostByID(ctx context.Context, id int) (*entity.Post, error)
	UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error)
	DeletePost(ctx context.Context, id int) error
}

type postUsecase struct {
	postRepo repository.PostRepository
}

func NewPostUsecases(postRepo repository.PostRepository) PostUsecase {
	return &postUsecase{postRepo: postRepo}
}

func (u *postUsecase) CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	createdPost, err := u.postRepo.CreatePost(ctx, post)
	if err != nil {
		return nil, err
	}
	return createdPost, nil
}

func (u *postUsecase) GetPosts(ctx context.Context) ([]entity.Post, error) {
	posts, err := u.postRepo.GetPosts(ctx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (u *postUsecase) GetPostByID(ctx context.Context, id int) (*entity.Post, error) {
	post, err := u.postRepo.GetPostByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (u *postUsecase) UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	updatedPost, err := u.postRepo.UpdatePost(ctx, post)
	if err != nil {
		return nil, err
	}
	return updatedPost, nil
}

func (u *postUsecase) DeletePost(ctx context.Context, id int) error {
	err := u.postRepo.DeletePost(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
