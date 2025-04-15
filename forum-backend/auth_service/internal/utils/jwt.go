package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTUtil struct {
	secret string
}

func NewJWTUtil(secret string) *JWTUtil {
	return &JWTUtil{secret: secret}
}

func (j *JWTUtil) GenerateToken(userID int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(j.secret))
}

func (j *JWTUtil) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
}
