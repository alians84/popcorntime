package main

import (
	"log"
	_ "popcorntime-project/docs"
	"popcorntime-project/internal/config"
	"popcorntime-project/internal/routes"
	"time"
)

// @title PopcornTime API
// @version 1.0
// @description API для синхронизированного просмотра видео
// @host popcorntimes.ru
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config.ConnectPostgres()
	defer func() {
		if sqlDB, err := config.DB.DB(); err == nil {
			sqlDB.Close() // Закрытие соединения при завершении
		}
	}()

	// Конфигурация Fiber
	fiberCfg := config.FiberConfig{
		Prefork:      false, // Для разработки (в продакшене — true)
		ServerHeader: "PopcornTime API",
		IdleTimeout:  30 * time.Second, // Для WebSocket-подключений
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app := config.SetupFiber(fiberCfg)

	routes.SetupRoutes(app)

	// Запуск сервера
	log.Fatal(app.Listen(":3000"))
}
