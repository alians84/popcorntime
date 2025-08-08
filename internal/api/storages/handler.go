package storages

import (
	"github.com/gofiber/fiber/v2"
	_ "popcorntime-project/internal/models"
	"popcorntime-project/internal/repository/storages"
)

type Handler struct {
	service storages.Service
}

func NewHandler(service storages.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Storages godoc
// @Summary Загрузить изображение в B2 Cloud Storage
// @Description Загружает изображение в Backblaze B2 Cloud Storage
// @Tags Storages
// @Accept multipart/form-data
// @Param file formData file true "Файл изображения"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /api/storages/upload [post]
func (h *Handler) UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Не удалось получить файл",
		})
	}

	uploadFile, err := h.service.UploadFile(c.Context(), file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error upload file to s3",
		})
	}

	return c.JSON(
		fiber.Map{
			"message":    "success",
			"uploadFile": uploadFile,
		},
	)
}
