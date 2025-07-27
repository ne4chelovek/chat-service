package chatHandler

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/ne4chelovek/chat_service/internal/model"
	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"strconv"
	"sync"
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
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer ws.Close()

	// 1. Читаем начальное сообщение с токеном
	var tokenMsg struct {
		Token string `json:"token"`
	}

	if err := ws.ReadJSON(&tokenMsg); err != nil {
		log.Printf("Failed to read token message: %v", err)
		return
	}

	if tokenMsg.Token == "" {
		log.Printf("Empty token")
		ws.WriteJSON(map[string]string{"error": "token required"})
		return
	}

	ctx := metadata.AppendToOutgoingContext(r.Context(), "authorization", tokenMsg.Token)

	var initMsg model.WsMessage

	if err := ws.ReadJSON(&initMsg); err != nil {
		log.Printf("Failed to read initial message: %v", err)
		return
	}

	chatID := strconv.Itoa(int(initMsg.ChatID))

	// 2. Подключаемся к gRPC stream
	stream, err := c.chatClient.ConnectChat(ctx, &desc.ConnectChatRequest{
		ChatId:   chatID,
		Username: initMsg.From,
	})
	if err != nil {
		log.Printf("gRPC stream error: %v", err)
		ws.WriteJSON(map[string]string{"error": "failed to connect to chat"})
		return
	}

	// 3. для управления горутинами
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 4. Горутина для чтения из WebSocket и отправки в gRPC
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
			}
		}()

		for {
			var msg model.WsMessage
			if err := ws.ReadJSON(&msg); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Printf("WebSocket read error: %v", err)
				}
				wg.Done()
				return
			}

			// Отправляем сообщение через gRPC
			_, err = c.chatClient.SendMessage(ctx, &desc.SendMessageRequest{
				ChatId: msg.ChatID,
				Message: &desc.Message{
					Text: msg.Text,
					From: msg.From,
				},
			})
			if err != nil {
				log.Printf("gRPC SendMessage failed: %v", err)
				wg.Done()
				return
			}
		}
	}()

	// 5. Горутина для чтения из gRPC stream и отправки в WebSocket
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := stream.Recv()
				if err != nil {
					log.Printf("gRPC stream receive error: %v", err)
					wg.Done()
					return
				}

				if err := ws.WriteJSON(msg); err != nil {
					log.Printf("WebSocket write error: %v", err)
					wg.Done()
					return
				}
			}
		}
	}()

	// 6. Ожидаем завершения одной из горутин
	wg.Wait()
}
