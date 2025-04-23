package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/internal/repository/adapters"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestPostRepository_CreatePost_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание адаптера
	dbAdapter := adapters.DbAdapter{db}

	// Инициализация репозитория
	postRepo := NewPostRepository(&dbAdapter, logger)

	// Тестовые данные
	post := entity.Post{
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}
	createdPost := post
	createdPost.ID = 1

	// Настройка моков
	mock.ExpectExec(`INSERT INTO posts \(author_id, title, content\) VALUES \(\?, \?, \?\)`).
		WithArgs(post.AuthorId, post.Title, post.Content).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызов метода
	result, err := postRepo.CreatePost(context.Background(), post)

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, &createdPost, result)

	// Проверка вызовов моков
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_CreatePost_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание адаптера
	dbAdapter := &adapters.DbAdapter{DB: db}

	// Инициализация репозитория
	postRepo := NewPostRepository(dbAdapter, logger)

	// Тестовые данные
	post := entity.Post{
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}

	// Настройка моков
	mock.ExpectExec(`INSERT INTO posts \(author_id, title, content\) VALUES \(\?, \?, \?\)`).
		WithArgs(post.AuthorId, post.Title, post.Content).
		WillReturnError(errors.New("failed to create post"))

	// Вызов метода
	result, err := postRepo.CreatePost(context.Background(), post)

	// Проверка результата
	assert.Error(t, err)
	assert.Nil(t, result)

	// Проверка вызовов моков
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetPosts_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание адаптера
	dbAdapter := &adapters.DbAdapter{DB: db}

	// Инициализация репозитория
	postRepo := NewPostRepository(dbAdapter, logger)

	// Тестовые данные
	posts := []entity.Post{
		{ID: 1, AuthorId: 1, Title: "Post 1", Content: "Content 1"},
		{ID: 2, AuthorId: 2, Title: "Post 2", Content: "Content 2"},
	}

	// Настройка моков
	rows := sqlmock.NewRows([]string{"id", "author_id", "title", "content"})
	for _, post := range posts {
		rows.AddRow(post.ID, post.AuthorId, post.Title, post.Content)
	}
	mock.ExpectQuery(`SELECT id, author_id, title, content FROM posts`).
		WillReturnRows(rows)

	// Вызов метода
	result, err := postRepo.GetPosts(context.Background())

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, posts, result)

	// Проверка вызовов моков
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetPosts_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание адаптера
	dbAdapter := &adapters.DbAdapter{DB: db}

	// Инициализация репозитория
	postRepo := NewPostRepository(dbAdapter, logger)

	// Настройка моков
	mock.ExpectQuery(`SELECT id, author_id, title, content FROM posts`).
		WillReturnError(errors.New("failed to get posts"))

	// Вызов метода
	result, err := postRepo.GetPosts(context.Background())

	// Проверка результата
	assert.Error(t, err)
	assert.Nil(t, result)

	// Проверка вызовов моков
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetPostByID_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание адаптера
	dbAdapter := &adapters.DbAdapter{DB: db}

	// Инициализация репозитория
	postRepo := NewPostRepository(dbAdapter, logger)

	// Тестовые данные
	postID := 1

	// Настройка моков
	mock.ExpectQuery(`SELECT id, author_id, title, content FROM posts WHERE id = \?`).
		WithArgs(postID).
		WillReturnError(errors.New("failed to get post"))

	// Вызов метода
	result, err := postRepo.GetPostByID(context.Background(), postID)

	// Проверка результата
	assert.Error(t, err)
	assert.Nil(t, result)

	// Проверка вызовов моков
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_UpdatePost_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание адаптера
	dbAdapter := &adapters.DbAdapter{DB: db}

	// Инициализация репозитория
	postRepo := NewPostRepository(dbAdapter, logger)

	// Тестовые данные
	post := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Updated Post",
		Content:  "This is an updated post",
	}

	// Настройка моков
	mock.ExpectExec(`UPDATE posts SET title = \?, content = \? WHERE id = \?`).
		WithArgs(post.Title, post.Content, post.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызов метода
	result, err := postRepo.UpdatePost(context.Background(), post)

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, &post, result)

	// Проверка вызовов моков
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_UpdatePost_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание адаптера
	dbAdapter := &adapters.DbAdapter{DB: db}

	// Инициализация репозитория
	postRepo := NewPostRepository(dbAdapter, logger)

	// Тестовые данные
	post := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Updated Post",
		Content:  "This is an updated post",
	}

	// Настройка моков
	mock.ExpectExec(`UPDATE posts SET title = \?, content = \? WHERE id = \?`).
		WithArgs(post.Title, post.Content, post.ID).
		WillReturnError(errors.New("failed to update post"))

	// Вызов метода
	result, err := postRepo.UpdatePost(context.Background(), post)

	// Проверка результата
	assert.Error(t, err)
	assert.Nil(t, result)

	// Проверка вызовов моков
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_DeletePost_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание адаптера
	dbAdapter := &adapters.DbAdapter{DB: db}

	// Инициализация репозитория
	postRepo := NewPostRepository(dbAdapter, logger)

	// Тестовые данные
	postID := 1

	// Настройка моков
	mock.ExpectExec(`DELETE FROM posts WHERE id = \?`).
		WithArgs(postID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызов метода
	err = postRepo.DeletePost(context.Background(), postID)

	// Проверка результата
	assert.NoError(t, err)

	// Проверка вызовов моков
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_DeletePost_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание адаптера
	dbAdapter := &adapters.DbAdapter{DB: db}

	// Инициализация репозитория
	postRepo := NewPostRepository(dbAdapter, logger)

	// Тестовые данные
	postID := 1

	// Настройка моков
	mock.ExpectExec(`DELETE FROM posts WHERE id = \?`).
		WithArgs(postID).
		WillReturnError(errors.New("failed to delete post"))

	// Вызов метода
	err = postRepo.DeletePost(context.Background(), postID)

	// Проверка результата
	assert.Error(t, err)

	// Проверка вызовов моков
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetUserIDByToken_Success(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание адаптера
	dbAdapter := &adapters.DbAdapter{DB: db}

	// Инициализация репозитория
	postRepo := NewPostRepository(dbAdapter, logger)

	// Тестовые данные
	token := "test-token"
	userID := 1

	// Настройка моков
	rows := sqlmock.NewRows([]string{"user_id"}).AddRow(userID)
	mock.ExpectQuery(`SELECT user_id FROM tokens WHERE token = \?`).
		WithArgs(token).
		WillReturnRows(rows)

	// Вызов метода
	result, err := postRepo.GetUserIDByToken(context.Background(), token)

	// Проверка результата
	assert.NoError(t, err)
	assert.Equal(t, userID, result)

	// Проверка вызовов моков
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetUserIDByToken_Failure(t *testing.T) {
	// Инициализация логгера
	logger, _ := zap.NewProduction()

	// Создание моков
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Создание адаптера
	dbAdapter := &adapters.DbAdapter{DB: db}

	// Инициализация репозитория
	postRepo := NewPostRepository(dbAdapter, logger)

	// Тестовые данные
	token := "test-token"

	// Настройка моков
	mock.ExpectQuery(`SELECT user_id FROM tokens WHERE token = \?`).
		WithArgs(token).
		WillReturnError(errors.New("failed to get user ID by token"))

	// Вызов метода
	result, err := postRepo.GetUserIDByToken(context.Background(), token)

	// Проверка результата
	assert.Error(t, err)
	assert.Equal(t, 0, result)

	// Проверка вызовов моков
	assert.NoError(t, mock.ExpectationsWereMet())
}
