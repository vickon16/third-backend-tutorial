package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vickon16/third-backend-tutorial/cmd/types"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func CreateJWT(userId, userEmail string) (string, error) {
	expiration := time.Now().Add(time.Duration(Configs.JWTExpirationInSeconds) * time.Second)

	claims := types.JWTClaims{
		UserId:    userId,
		UserEmail: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(Configs.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
