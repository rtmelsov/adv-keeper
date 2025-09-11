// Package helpers
package helpers

// jwt.go
import (
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

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAccessJWT(userID string) (string, time.Time, error) {
	envs, err := LoadConfig()
	if err != nil {
		return "", time.Now(), err
	}

	now := time.Now().UTC()
	exp := now.Add(mustParseDuration(envs.AccessTTL))

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,                         // дублируем id в стандартный sub
			Audience:  jwt.ClaimStrings{"adv-keeper"}, // aud
			IssuedAt:  jwt.NewNumericDate(now),        // iat
			ExpiresAt: jwt.NewNumericDate(exp),        // exp
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(envs.JWTSecret)) // ключ как []byte
	return s, exp, err
}

func VerifyToken(tokenStr string) (*Claims, error) {
	envs, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	// Парсер сразу проверит метод подписи и audience
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithAudience("adv-keeper"),
	)

	claims := &Claims{}
	token, err := parser.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return []byte(envs.JWTSecret), nil // ключ как []byte
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Обычно парсер уже проверяет exp/nbf/iat; на всякий случай:
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now().UTC()) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
