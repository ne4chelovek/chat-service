package chat

import (
	"context"

	"fmt"
)

func (s *serv) GetChatInfo(ctx context.Context, chatID int64) ([]string, error) {
	var users []string
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		users, errTx = s.chatRepository.GetChatInfo(ctx, chatID)
		if errTx != nil {
			return errTx
		}

		log := fmt.Sprintf("get chat with id: %v", chatID)
		errTx = s.logRepository.Log(ctx, log)
		if errTx != nil {
			fmt.Printf("Log error: %v\n", errTx)
			return errTx
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Transacrion error: %v\n", err)
		return nil, err
	}
	return users, nil
}
