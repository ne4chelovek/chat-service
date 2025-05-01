package chat

import (
	"github.com/ne4chelovek/chat_common/pkg/db"
	"github.com/ne4chelovek/chat_service/internal/repository"
	"github.com/ne4chelovek/chat_service/internal/service"
)

type serv struct {
	chatRepository repository.ChatRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager
}

func NewService(chatRepository repository.ChatRepository, logRepository repository.LogRepository, txManager db.TxManager) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		logRepository:  logRepository,
		txManager:      txManager,
	}
}
