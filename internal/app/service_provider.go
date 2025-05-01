package app

import (
	"context"
	"github.com/ne4chelovek/chat_common/pkg/closer"
	"github.com/ne4chelovek/chat_common/pkg/db"
	"github.com/ne4chelovek/chat_common/pkg/db/pg"
	"github.com/ne4chelovek/chat_common/pkg/db/transaction"
	"github.com/ne4chelovek/chat_service/internal/api/chat"
	"github.com/ne4chelovek/chat_service/internal/repository"
	"github.com/ne4chelovek/chat_service/internal/service"
	"log"

	chatRepository "github.com/ne4chelovek/chat_service/internal/repository/chat"
	logRepository "github.com/ne4chelovek/chat_service/internal/repository/log"
	chatService "github.com/ne4chelovek/chat_service/internal/service/chat"
)

const (
	dbDSN = "host=localhost port=54321 dbname=chat user=chat-user password=chat-password sslmode=disable"
)

type serviceProvider struct {
	dbClient  db.Client
	txManager db.TxManager

	chatRepository repository.ChatRepository
	chatService    service.ChatService
	logRepository  repository.LogRepository

	chatServ *chat.Server
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, dbDSN)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}
	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
	}
	return s.chatRepository
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}
	return s.logRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.ChatRepository(ctx),
			s.LogRepository(ctx),
			s.TxManager(ctx),
		)
	}
	return s.chatService
}

func (s *serviceProvider) ChatImpl(ctx context.Context) *chat.Server {
	if s.chatServ == nil {
		s.chatServ = chat.NewImplementation(s.ChatService(ctx))
	}
	return s.chatServ
}
