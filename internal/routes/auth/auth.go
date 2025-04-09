package auth

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"popcorntime-project/internal/api/auth"
	"popcorntime-project/internal/middleware"
)

func SetupAuthRoutes(app *fiber.App, handler *auth.AuthHandler) {
	authGroup := app.Group("api/auth")
	{
		authGroup.Post("/register", handler.Register)
		authGroup.Post("/login", handler.Login)
		authGroup.Post("/refresh", handler.Refresh)

		authGroup.Get("/profile", middleware.JWTProtected(os.Getenv("JWT_SECRET")), func(c *fiber.Ctx) error {
			return nil
		})
	}
}
