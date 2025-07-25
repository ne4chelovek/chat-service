package chat

import (
	"github.com/ne4chelovek/chat_common/pkg/db"
	"github.com/ne4chelovek/chat_service/internal/model"
	catApi "github.com/ne4chelovek/chat_service/internal/openApi"
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
	api            catApi.ApiCat

	channels   map[string]chan *model.Message
	mxChannels sync.RWMutex

	chats   map[string]*chat
	mxChats sync.RWMutex
}

type chat struct {
	streams map[string]model.Stream
	m       sync.RWMutex
}

func NewService(chatRepository repository.ChatRepository, logRepository repository.LogRepository, api catApi.ApiCat, txManager db.TxManager) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		logRepository:  logRepository,
		txManager:      txManager,
		api:            api,
		channels:       make(map[string]chan *model.Message),
		chats:          make(map[string]*chat),
	}
}
