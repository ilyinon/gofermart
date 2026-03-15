package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("super-secret")

func GenerateToken(userID int64) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		},
	)

	return token.SignedString(secret)
}

func ParseToken(tokenString string) (int64, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	)

	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)

	id := int64(claims["user_id"].(float64))

	return id, nil
}
