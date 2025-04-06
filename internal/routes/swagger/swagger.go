package swagger

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func SetupAPIRoutes(app *fiber.App) {
	// Swagger
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Health check
	app.Get("/api/health", healthCheck)

	app.Get("api/check", healthCheckTwo)
}

// @Summary Проверка здоровья сервера
// @Success 200 {object} map[string]string
// @Router /api/health [get]
func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "server jobs"})
}

// @Summary Проверка здоровья сервера
// @Success 200 {object} map[string]string
// @Router /api/health [get]
func healthCheckTwo(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "server jobs"})
}
