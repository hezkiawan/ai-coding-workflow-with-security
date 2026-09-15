package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"opsdesk/internal/models"
)

type Claims struct {
	UserID     int64       `json:"user_id"`
	Username   string      `json:"username"`
	Email      string      `json:"email"`
	Role       models.Role `json:"role"`
	TokenID    string      `json:"jti,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(user *models.User, secret string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "opsdesk-auth",
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenStr string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
