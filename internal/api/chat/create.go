package chat

import (
	"context"
	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"
)

func (s *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.chatService.Create(ctx, req.GetUsernames())
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{
		Id: id,
	}, nil

}
