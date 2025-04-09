package users

import (
	"gorm.io/gorm"
	"popcorntime-project/internal/repository/users"
)

func Init(db *gorm.DB) *Handler {
	// Инициализация сервиса
	userService := users.NewService(db)

	// Создание обработчика с внедренным сервисом
	return NewHandler(userService)
}
