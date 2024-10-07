package token

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenUseCase interface {
	GenerateAccessToken(claims JwtCustomClaims) (string, error)
}

type tokenUseCase struct {
	secretKey string
}

type JwtCustomClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenUseCase(secretKey string) *tokenUseCase {
	return &tokenUseCase{
		secretKey: secretKey,
	}
}

func (t *tokenUseCase) GenerateAccessToken(claims JwtCustomClaims) (string, error) {
	plainToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	encodedToken, err := plainToken.SignedString([]byte(t.secretKey))

	if err != nil {
		return "", err
	}

	return encodedToken, nil
}
