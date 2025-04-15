package repository

import (
	"forum/auth_service/internal/entity"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) Register(user entity.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", user.Username, user.Password, user.Role)
	return err
}

func (r *AuthRepository) GetUserByUsername(username string) (entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, "SELECT id, username, password, role FROM users WHERE username=?", username)
	return user, err
}

func (r *AuthRepository) SaveToken(userID int, token string) error {
	_, err := r.db.Exec("INSERT INTO tokens (user_id, token) VALUES (?, ?)", userID, token)
	return err
}
