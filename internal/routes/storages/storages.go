package storages

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"popcorntime-project/internal/api/storages"
	"popcorntime-project/internal/middleware"
)

func SetupStorageRoutes(app *fiber.App, handler *storages.Handler) {
	storageGroup := app.Group("api/storages")
	{
		storageGroup.Use(middleware.JWTProtected(os.Getenv("JWT_SECRET")))

		storageGroup.Post("/upload", handler.UploadImage)

	}
	// Роуты для загрузки изображений
	//upload := api.Group("/upload")
	//upload.Get("/url", handler.GenerateUploadURL)
	//upload.Post("/complete", handler.CompleteUpload)
}
