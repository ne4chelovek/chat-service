package chat

import (
	"context"
	"fmt"
	"github.com/ne4chelovek/chat_service/internal/model"
)

func (s *serv) SendMessage(ctx context.Context, chatID int64, mes *model.Message) (string, error) {
	var status string
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		status, errTx = s.chatRepository.SendMessage(ctx, chatID, mes)
		if errTx != nil {
			return errTx
		}
		log := fmt.Sprintf("SendMessage in chat with id: %d, message: %v", chatID, mes)
		errTx = s.logRepository.Log(ctx, log)
		if errTx != nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return status, nil
}
