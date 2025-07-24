package service

import (
	"context"
	"github.com/ne4chelovek/chat_service/internal/model"
	"net/http"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatService interface {
	Create(ctx context.Context, user []string) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) (*emptypb.Empty, error)
	GetChatInfo(ctx context.Context, chatID int64) ([]string, error)
	SendMessage(ctx context.Context, chatID int64, mes *model.Message) (string, error)
	Connect(chatID int64, username string, stream model.Stream) error
}

type HandlerService interface {
	HandleChatWebSocket(w http.ResponseWriter, r *http.Request)
}
