package auth

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"popcorntime-project/internal/models"
	"popcorntime-project/internal/repository/users"
	"popcorntime-project/pkg/utils"
)

type AuthHandler struct {
	service   users.Service
	jwtSecret string
}

func NewHandler(service users.Service, jwt string) *AuthHandler {
	return &AuthHandler{
		service:   service,
		jwtSecret: jwt,
	}
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Вход в систему по email и паролю
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Данные для входа"
// @Success 200 {object} models.TokenResponse "Успешная аутентификация"
// @Failure 400 {object} models.ErrorResponse "Неверный формат запроса"
// @Failure 401 {object} models.ErrorResponse "Неверные учетные данные"
// @Failure 500 {object} models.ErrorResponse "Ошибка сервера"
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Проверка пользователя
	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Генерация токенов
	accessToken, err := utils.GenerateAccessToken(uint(user.ID), h.jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate refresh token"})
	}

	return c.JSON(models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(utils.AccessTokenExpire.Seconds()),
		User: models.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	})
}

// Refresh godoc
// @Summary Обновление токенов
// @Description Получение новой пары access/refresh токенов
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RefreshRequest true "Refresh token"
// @Success 200 {object} models.TokenResponse
// @Failure 401 {object} map[string]string
// @Router /api/auth/refresh [post]
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	// 1. Проверяем access token (даже просроченный)
	tokenString := utils.ExtractToken(c)
	claims, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.MapClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(h.jwtSecret), nil
		},
	)

	// Разрешаем запрос, даже если токен просрочен
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// 2. Извлекаем user_id
	userID := (*claims.Claims.(*jwt.MapClaims))["user_id"].(float64)

	// 3. Проверяем refresh token
	var req struct{ RefreshToken string }
	if err := c.BodyParser(&req); err != nil || req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Refresh token required",
		})
	}

	// 4. Генерируем новую пару
	return h.generateTokenResponse(c, uint(userID))
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создание нового аккаунта
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "Данные регистрации"
// @Success 201 {object} models.TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}

	// Валидация (можно использовать github.com/go-ozzo/ozzo-validation)
	if req.Email == "" || len(req.Password) < 8 || req.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email, пароль (мин. 8 символов) и имя пользователя обязательны",
		})
	}

	// Создание пользователя
	user, err := h.service.Register(req.Email, req.Password, req.Username)
	if err != nil {
		//if utils.IsDuplicateKeyError(err) {
		//	return c.Status(fiber.StatusConflict).JSON(fiber.Map{
		//		"error": "Email или имя пользователя уже заняты",
		//	})
		//}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при создании пользователя",
		})
	}

	// Генерация токенов (как в логине)
	accessToken, err := utils.GenerateAccessToken(uint(user.ID), h.jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка генерации токена",
		})
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка генерации refresh-токена",
		})
	}

	return c.JSON(models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(utils.AccessTokenExpire.Seconds()),
		User: models.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	})
}

// generateTokenResponse создает новую пару токенов и возвращает ответ
func (h *AuthHandler) generateTokenResponse(c *fiber.Ctx, userID uint) error {
	// 1. Генерация access token
	accessToken, err := utils.GenerateAccessToken(userID, h.jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate access token",
		})
	}

	// 2. Генерация refresh token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	// 3. Получаем данные пользователя для ответа
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user data",
		})
	}

	// 4. Формируем ответ
	return c.JSON(models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(utils.AccessTokenExpire.Seconds()),
		User: models.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	})
}
