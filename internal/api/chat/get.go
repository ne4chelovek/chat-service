package chat

import (
	"context"
	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"
	"log"
)

func (s *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	users, err := s.chatService.GetChatInfo(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("get chat with id: %d", req.GetId())

	return &desc.GetResponse{
		Usernames: users,
	}, nil
}
