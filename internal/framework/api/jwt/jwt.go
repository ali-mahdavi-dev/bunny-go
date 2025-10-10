package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(exp time.Duration, jwtSecret string, user_id uint64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user_id,
		"exp":     time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
