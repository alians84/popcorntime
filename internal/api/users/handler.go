package users

import (
	"github.com/gofiber/fiber/v2"
	_ "popcorntime-project/internal/models"
	"popcorntime-project/internal/repository/users"
	"strconv"
)

type Handler struct {
	service users.Service
}

func NewHandler(service users.Service) *Handler {
	return &Handler{service: service}
}

// GetUser godoc
// @Summary Получить пользователя
// @Description Получить информацию о пользователе по ID
// @Tags Users
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @scheme Bearer
// @Router /api/users/{id} [get]
func (h *Handler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.service.GetUser(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

// GetUsers godoc
// @Summary Получить пользователей
// @Description Получить информацию о пользователях
// @Tags Users
// @Produce json
// @Success 200 {object} []models.User
// @Failure 404 {object} map[string]string
// @Router /api/users [get]
func (h *Handler) GetUsers(c *fiber.Ctx) error {
	getUsers, err := h.service.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Users not found",
		})
	}

	return c.JSON(getUsers)
}
