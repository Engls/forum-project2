package usecase

import (
	"errors"
	"forum/auth_service/internal/entity"
	"forum/auth_service/internal/repository"
	"forum/auth_service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	authRepo *repository.AuthRepository
	jwtUtil  *utils.JWTUtil
}

func NewAuthUsecase(authRepo *repository.AuthRepository, jwtUtil *utils.JWTUtil) *AuthUsecase {
	return &AuthUsecase{authRepo: authRepo, jwtUtil: jwtUtil}
}

func (u *AuthUsecase) Register(username, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := entity.User{Username: username, Password: string(hashedPassword), Role: role}
	return u.authRepo.Register(user)
}

func (u *AuthUsecase) Login(username, password string) (string, error) {
	user, err := u.authRepo.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	token, err := u.jwtUtil.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}
	if err := u.authRepo.SaveToken(user.ID, token); err != nil {
		return "", err
	}
	return token, nil
}
