package routes

import (
	"github.com/gofiber/fiber/v2"
	initauth "popcorntime-project/internal/api/auth"
	initstorages "popcorntime-project/internal/api/storages"
	inituser "popcorntime-project/internal/api/users"
	"popcorntime-project/internal/config"
	"popcorntime-project/internal/routes/auth"
	"popcorntime-project/internal/routes/storages"
	"popcorntime-project/internal/routes/swagger"
	"popcorntime-project/internal/routes/users"
	"popcorntime-project/internal/routes/ws"
)

func SetupRoutes(app *fiber.App, b2Client *config.B2Client) {
	users.SetupUserRoutes(app, inituser.Init(config.DB))
	swagger.SetupAPIRoutes(app)
	ws.SetupWSRoutes(app)
	auth.SetupAuthRoutes(app, initauth.Init(config.DB))
	storages.SetupStorageRoutes(app, initstorages.Init(config.DB, b2Client))
}
