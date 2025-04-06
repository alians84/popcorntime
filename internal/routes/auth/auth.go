package auth

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func SetupAuthRoutes(app *fiber.App) {
	authGroup := app.Group("/auth")

	authGroup.Get("/:provider", oauthLogin)
	authGroup.Get("/:provider/callback", oauthCallback)
}

// @Summary Вход через OAuth провайдера
// @Param provider path string true "google/github"
// @Router /auth/{provider} [get]
func oauthLogin(c *fiber.Ctx) error {
	// Логика OAuth
	return http.ErrAbortHandler
}

// @Summary Callback от OAuth провайдера
// @Router /auth/{provider}/callback [get]
func oauthCallback(c *fiber.Ctx) error {
	// Логика обработки callback
	return http.ErrAbortHandler
}
