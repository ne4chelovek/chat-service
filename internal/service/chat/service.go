package chat

import (
	"github.com/ne4chelovek/chat_common/pkg/db"
	"github.com/ne4chelovek/chat_service/internal/model"
	"github.com/ne4chelovek/chat_service/internal/repository"
	"github.com/ne4chelovek/chat_service/internal/service"
	"sync"
)

const (
	messagesBuffer = 100
)

type serv struct {
	chatRepository repository.ChatRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager

	channels   map[string]chan *model.Message
	mxChannels sync.RWMutex

	chats   map[string]*chat
	mxChats sync.RWMutex
}

type chat struct {
	streams map[string]model.Stream
	m       sync.RWMutex
}

func NewService(chatRepository repository.ChatRepository, logRepository repository.LogRepository, txManager db.TxManager) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		logRepository:  logRepository,
		txManager:      txManager,
		channels:       make(map[string]chan *model.Message),
		chats:          make(map[string]*chat),
	}
}
