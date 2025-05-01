package chat

import (
	"github.com/ne4chelovek/chat_common/pkg/db"
	"github.com/ne4chelovek/chat_service/internal/model"
	"github.com/ne4chelovek/chat_service/internal/repository"

	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	tableChats    = "chats"
	tableMessages = "messages"

	idColumn         = "chat_id"
	usernamesColumn  = "usernames"
	from_userColumn  = "from_user"
	textColumn       = "text"
	created_atColumn = "created_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, users []string) (int64, error) {
	builderInsert := sq.Insert(tableChats).
		PlaceholderFormat(sq.Dollar).
		Columns(created_atColumn, usernamesColumn).
		Values(sq.Expr("NOW()"), users).
		Suffix("RETURNING chat_id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRaw: query,
	}

	var chat_id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chat_id)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}

	return chat_id, nil
}

func (r *repo) GetChatInfo(ctx context.Context, chatID int64) ([]string, error) {
	builderSelect := sq.Select(usernamesColumn).
		From(tableChats).
		Where(sq.Eq{idColumn: chatID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "chat_repository.GetChatInfo",
		QueryRaw: query,
	}

	var usernames []string // Используем срез строк вместо pgtype.TextArray
	err = r.db.DB().ScanOneContext(ctx, &usernames, q, args...)
	if err != nil {
		return nil, fmt.Errorf("scan error: %v", err)
	}

	return usernames, nil
}

func (r *repo) DeleteChat(ctx context.Context, chatID int64) (*emptypb.Empty, error) {
	builderDelete := sq.Delete(tableChats).
		Where(sq.Eq{idColumn: chatID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query %v:", err)
	}

	q := db.Query{
		Name:     "chat_repository.DeleteChat",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %v: ", err)
	}
	return nil, nil
}

func (r *repo) SendMessage(ctx context.Context, chatID int64, mes *model.Message) (string, error) {
	builderInsert := sq.Insert(tableMessages).
		PlaceholderFormat(sq.Dollar).
		Columns(idColumn, from_userColumn, textColumn).
		Values(chatID, mes.From, mes.Text).
		Suffix("RETURNING status")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return "", fmt.Errorf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "chat_repository.SendMessage",
		QueryRaw: query,
	}

	var status string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %v", err)
	}
	return status, nil
}
