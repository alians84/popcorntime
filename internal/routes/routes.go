package routes

import (
	"github.com/gofiber/fiber/v2"
	"popcorntime-project/internal/routes/auth"
	"popcorntime-project/internal/routes/swagger"
	"popcorntime-project/internal/routes/ws"
)

func SetupRoutes(app *fiber.App) {
	swagger.SetupAPIRoutes(app)
	ws.SetupWSRoutes(app)
	auth.SetupAuthRoutes(app)
}
