package chat

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serv) DeleteChat(ctx context.Context, chatID int64) (*emptypb.Empty, error) {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		_, errTx = s.chatRepository.DeleteChat(ctx, chatID)
		if errTx != nil {
			return errTx
		}
		log := fmt.Sprintf("deleted chat with id: %v", chatID)
		errTx = s.logRepository.Log(ctx, log)
		if errTx != nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
