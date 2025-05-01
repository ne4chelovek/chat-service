package chat

import (
	"github.com/ne4chelovek/chat_service/internal/service"
	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"
)

type Server struct {
	desc.UnimplementedChatServer
	chatService service.ChatService
}

func NewImplementation(chatService service.ChatService) *Server {
	return &Server{
		chatService: chatService,
	}
}
