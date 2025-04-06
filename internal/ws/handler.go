package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

func HandleConnections(c *websocket.Conn) {
	defer c.Close()

	for {
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			log.Error("WebSocket read error:", err)
			break
		}

		// Пример: Обработка событий (play, pause, chat)
		switch string(msg) {
		case "play":
			broadcastToRoom(c, "play")
		case "pause":
			broadcastToRoom(c, "pause")
		default:
			log.Info("Received:", string(msg))
		}

		if err := c.WriteMessage(msgType, msg); err != nil {
			log.Error("WebSocket write error:", err)
			break
		}
	}
}

func broadcastToRoom(conn *websocket.Conn, message string) {
	// Логика рассылки сообщений в комнату
}
