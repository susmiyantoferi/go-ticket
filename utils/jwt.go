package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaim struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	Name         string `json:"name"`
	Token        string `json:"access_token"`
	TokenRefresh string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExipresIn    int    `json:"expires_in"`
}

func GenerateToken(userId uint, name, email, role string, exp time.Duration) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	jwtExp := time.Now().Add(exp * time.Hour)

	tokenCLaim := &TokenClaim{
		UserID: userId,
		Name:   name,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(jwtExp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenCLaim)
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenStr, nil
}

func ClaimTokenRefresh(tokenUser string) (*TokenClaim, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	claim := &TokenClaim{}

	token, err := jwt.ParseWithClaims(tokenUser, claim, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")

	}

	return claim, nil
}
