package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret []byte

func MustInitJWT(s string) {
	if s == "" {
		panic("PANIC :: JwtSecret is not set.")
	}
	JwtSecret = []byte(s)
}

type AccessTokenClaims struct {
	UserID  uint   `json:"user_id"`
	TokenID string `json:"jti"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uint) (string, string, error) {
	tokenID := generateTokenID()

	claims := AccessTokenClaims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(JwtSecret)
	return signed, tokenID, err
}

func ValidateAccessToken(tokenStr string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

func GenerateRefreshToken() string {
	b := make([]byte, 64)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateTokenID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
