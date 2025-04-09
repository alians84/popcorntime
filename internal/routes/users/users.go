package users

import (
	"github.com/gofiber/fiber/v2"
	"popcorntime-project/internal/api/users"
)

func SetupUserRoutes(router fiber.Router, handler *users.Handler) {
	userGroup := router.Group("api/users")
	userGroup.Get("/:id", handler.GetUser)
}
