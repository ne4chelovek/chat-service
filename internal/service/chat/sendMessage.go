package chat

import (
	"context"
	"github.com/ne4chelovek/chat_service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

func (s *serv) SendMessage(ctx context.Context, chatID int64, mes *model.Message) (string, error) {
	id := strconv.Itoa(int(chatID))
	s.mxChannels.RLock()
	chatChan, ok := s.channels[id]
	s.mxChannels.RUnlock()
	if !ok {
		return "", status.Errorf(codes.NotFound, "chat not found")
	}

	sts, err := s.chatRepository.SendMessage(ctx, chatID, mes)
	if err != nil {
		return "", err
	}

	chatChan <- mes

	return sts, nil
}
