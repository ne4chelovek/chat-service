package chat

import (
	"context"
	"github.com/ne4chelovek/chat_service/internal/converter"
	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"
)

func (s *Server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*desc.SendMessageResponse, error) {
	status, err := s.chatService.SendMessage(ctx, req.GetChatId(), converter.ToChatFromDesc(req.GetMessage()))
	if err != nil {
		return nil, err
	}
	return &desc.SendMessageResponse{
		Status: status,
	}, nil
}
