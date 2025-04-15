package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTUtil struct {
	secretKey string
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func NewJWTUtil(secretKey string) *JWTUtil {
	return &JWTUtil{secretKey: secretKey}
}

func (j *JWTUtil) GenerateToken(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
