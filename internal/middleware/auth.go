package middleware

import (
	"github.com/gofiber/fiber/v2"
	"popcorntime-project/pkg/utils"
)

func JWTProtected(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := extractToken(c)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header missing",
			})
		}

		userID, err := utils.ParseAccessToken(tokenString, secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		c.Locals("userID", userID)
		return c.Next()
	}
}

func extractToken(c *fiber.Ctx) string {
	bearer := c.Get("Authorization")
	if len(bearer) > 7 && bearer[:7] == "Bearer " {
		return bearer[7:]
	}
	return ""
}
