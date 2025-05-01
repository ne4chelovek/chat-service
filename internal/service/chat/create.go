package chat

import (
	"context"
	"fmt"
)

func (s *serv) Create(ctx context.Context, req []string) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		id, errTx = s.chatRepository.Create(ctx, req)
		if errTx != nil {
			return errTx
		}
		log := fmt.Sprintf("Create chat with id %v", id)
		errTx = s.logRepository.Log(ctx, log)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
