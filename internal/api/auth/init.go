package auth

import (
	"gorm.io/gorm"
	"os"
	"popcorntime-project/internal/repository/users"
)

func Init(db *gorm.DB) *AuthHandler {
	// Инициализация сервиса
	userService := users.NewService(db)

	// Создание обработчика с внедренным сервисом
	return NewHandler(
		userService,
		os.Getenv("JWT_SECRET"),
	)
}
