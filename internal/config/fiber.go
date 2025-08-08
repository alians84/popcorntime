package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"time"
)

// FiberConfig — кастомные настройки сервера
type FiberConfig struct {
	Prefork      bool
	ServerHeader string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func SetupFiber(cfg FiberConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       cfg.Prefork,      // Для многопоточной нагрузки
		ServerHeader:  cfg.ServerHeader, // "PopcornTime API"
		IdleTimeout:   cfg.IdleTimeout,  // 30s (для WebSocket)
		ReadTimeout:   cfg.ReadTimeout,  // 10s
		WriteTimeout:  cfg.WriteTimeout, // 10s
		CaseSensitive: true,             // Чувствительность к регистру в URL
		StrictRouting: false,            // /path и /path/ — одно и то же
	})

	// Мидлвари (порядок важен!)
	app.Use(
		helmet.New(), // Защита headers
		cors.New(cors.Config{ // CORS для фронтенда
			AllowOrigins:     "http://localhost:3000, http://127.0.0.1:3000, http://localhost:5173", // Добавлены популярные адреса для разработки
			AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
			AllowCredentials: true,
			MaxAge:           86400,
		}),
		logger.New(logger.Config{ // Логирование запросов
			Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
			TimeFormat: "2006-01-02 15:04:05",
		}),
	)

	return app
}
