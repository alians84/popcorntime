package users

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"popcorntime-project/internal/api/users"
	"popcorntime-project/internal/middleware"
)

func SetupUserRoutes(router fiber.Router, handler *users.Handler) {
	userGroup := router.Group("api/users")
	userGroup.Use(middleware.JWTProtected(os.Getenv("JWT_SECRET")))

	userGroup.Get("/:id", handler.GetUser)
	userGroup.Get("/", handler.GetUsers)
}
