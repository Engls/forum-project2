package usecase

import (
	"errors"
	"github.com/Engls/forum-project2/auth_service/internal/entity"
	"github.com/Engls/forum-project2/auth_service/internal/repository"
	"github.com/Engls/forum-project2/auth_service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(username, password, role string) error
	Login(username, password string) (string, error)
}

type authUsecase struct {
	authRepo repository.AuthRepository
	jwtUtil  *utils.JWTUtil
}

func NewAuthUsecase(authRepo repository.AuthRepository, jwtUtil *utils.JWTUtil) AuthUsecase {
	return &authUsecase{authRepo: authRepo, jwtUtil: jwtUtil}
}

func (u *authUsecase) Register(username, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := entity.User{Username: username, Password: string(hashedPassword), Role: role}
	return u.authRepo.Register(user)
}

func (u *authUsecase) Login(username, password string) (string, error) {
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
