package chatHandler

import (
	"github.com/gorilla/websocket"
	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"strconv"
)

type chatClient struct {
	chatClient desc.ChatClient
}

func NewChatClient(client desc.ChatClient) *chatClient {
	return &chatClient{
		chatClient: client,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Разрешаем все origin
}

func (c *chatClient) HandleChatWebSocket(w http.ResponseWriter, r *http.Request) {
	// Апгрейд соединения до WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer ws.Close()

	// Парсим параметры из URL
	query := r.URL.Query()
	chatIDStr := query.Get("chat_id")
	username := query.Get("username")

	_, err = strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Printf("Invalid chat ID: %v", err)
		return
	}

	// Создаем контекст с метаданными (если нужно)
	ctx := metadata.AppendToOutgoingContext(r.Context(), "username", username)

	// Подключаемся к gRPC streaming-ручке
	stream, err := c.chatClient.ConnectChat(ctx, &desc.ConnectChatRequest{
		ChatId:   chatIDStr,
		Username: username,
	})
	if err != nil {
		log.Printf("gRPC stream error: %v", err)
		return
	}

	// Запускаем обработку сообщений
	for {
		// Получаем сообщение из gRPC стрима
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("Stream receive error: %v", err)
			break
		}

		// Отправляем сообщение в WebSocket
		if err := ws.WriteJSON(msg); err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}
