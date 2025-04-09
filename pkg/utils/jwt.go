package utils

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

var (
	AccessTokenExpire  = 15 * time.Minute
	RefreshTokenExpire = 7 * 24 * time.Hour
)

// Генерация JWT
func GenerateAccessToken(userID uint, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(AccessTokenExpire).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Генерация refresh-токена
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Парсинг JWT
func ParseAccessToken(tokenString, secret string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return uint(claims["user_id"].(float64)), nil
}

func ExtractToken(c *fiber.Ctx) string {
	// 1. Проверка заголовка Authorization
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		// Формат: "Bearer eyJhbGciOiJIUzI1NiIsIn..."
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// 2. Проверка cookie (альтернативный вариант)
	token := c.Cookies("access_token")
	if token != "" {
		return token
	}

	// 3. Проверка query-параметра
	token = c.Query("token")
	if token != "" {
		return token
	}

	return ""
}
