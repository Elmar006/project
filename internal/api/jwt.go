package api

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	PasswordHash string `json:"pwd_hash"`
	jwt.RegisteredClaims
}

func generateToken(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))
	hashStr := hex.EncodeToString(hash[:])

	claims := &Claims{
		PasswordHash: hashStr,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(password))
}
func validateToken(tokenString, currentPassword string) bool {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(currentPassword), nil
	})

	if err != nil || !token.Valid {
		return false
	}
	hash := sha256.Sum256([]byte(currentPassword))
	expectedHash := hex.EncodeToString(hash[:])

	return claims.PasswordHash == expectedHash
}
