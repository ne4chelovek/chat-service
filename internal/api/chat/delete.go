package chat

import (
	"context"
	desc "github.com/ne4chelovek/chat_service/pkg/chat_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	_, err := s.chatService.DeleteChat(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return nil, nil
}
