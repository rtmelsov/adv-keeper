// Package helpers
package helpers

// jwt.go
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func mustParseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	if d == 0 {
		d = 15 * time.Minute
	}
	return d
}

func NewAccessJWT(userID string) (string, time.Time, error) {
	envs, err := LoadConfig()
	if err != nil {
		return "", time.Now(), err
	}

	now := time.Now().UTC()
	exp := now.Add(mustParseDuration(envs.AccessTTL))
	claims := jwt.MapClaims{
		"sub": userID,
		"aud": "adv-keeper",
		"iat": now.Unix(),
		"exp": exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(envs.JWTSecret))
	return s, exp, err
}

func SHA256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func VerifyToken(tokenStr string) (*Claims, error) {
	envs, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (any, error) {
		return envs.JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Проверим срок действия
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
