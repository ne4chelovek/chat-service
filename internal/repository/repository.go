package repository

import (
	"context"
	"github.com/ne4chelovek/chat_service/internal/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ChatRepository interface {
	Create(ctx context.Context, user []string) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) (*emptypb.Empty, error)
	GetChatInfo(ctx context.Context, chatID int64) ([]string, error)
	SendMessage(ctx context.Context, chatID int64, mes *model.Message) (string, error)
	GetMessage(ctx context.Context, chatID int64, page uint64) ([]*model.Message, error)
}

type LogRepository interface {
	Log(ctx context.Context, log string) error
}
