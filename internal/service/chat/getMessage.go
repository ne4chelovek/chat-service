package chat

import (
	"context"
	"github.com/ne4chelovek/chat_service/internal/model"
)

func (s *serv) GetMessage(ctx context.Context, chatID int64, page uint64) ([]*model.Message, error) {
	listMessage, err := s.chatRepository.GetMessage(ctx, chatID, page)
	if len(listMessage) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return listMessage, nil
}
