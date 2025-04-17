package repository

import (
	"github.com/Engls/forum-project2/auth_service/internal/entity"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	Register(user entity.User) error
	GetUserByUsername(username string) (entity.User, error)
	SaveToken(userID int, token string) error
}

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Register(user entity.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", user.Username, user.Password, user.Role)
	return err
}

func (r *authRepository) GetUserByUsername(username string) (entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, "SELECT id, username, password, role FROM users WHERE username=?", username)
	return user, err
}

func (r *authRepository) SaveToken(userID int, token string) error {
	_, err := r.db.Exec("INSERT INTO tokens (user_id, token) VALUES (?, ?)", userID, token)
	return err
}
