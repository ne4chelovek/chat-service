package app

import (
	"context"
	"fmt"
	"github.com/ne4chelovek/auth-service/pkg/access_v1"
	"github.com/ne4chelovek/chat_common/pkg/closer"
	"github.com/ne4chelovek/chat_common/pkg/db"
	"github.com/ne4chelovek/chat_common/pkg/db/pg"
	"github.com/ne4chelovek/chat_common/pkg/db/transaction"
	"github.com/ne4chelovek/chat_service/internal/api/chat"
	"github.com/ne4chelovek/chat_service/internal/client/rpc"
	rpcAuth "github.com/ne4chelovek/chat_service/internal/client/rpc/auth"
	clientApi "github.com/ne4chelovek/chat_service/internal/clientApi"
	"github.com/ne4chelovek/chat_service/internal/clientApi/apiHttp"
	"github.com/ne4chelovek/chat_service/internal/kafkaConsumer/consumer"
	"github.com/ne4chelovek/chat_service/internal/kafkaConsumer/handler"
	"github.com/ne4chelovek/chat_service/internal/repository"
	chatRepository "github.com/ne4chelovek/chat_service/internal/repository/chat"
	logRepository "github.com/ne4chelovek/chat_service/internal/repository/log"
	"github.com/ne4chelovek/chat_service/internal/service"
	chatService "github.com/ne4chelovek/chat_service/internal/service/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net/http"
	"time"
)

var kafkaAddresses = []string{
	"localhost:9091", // Для доступа с хоста
	"localhost:9092",
	"localhost:9093",
}

const (
	topic         = "user_session_events"
	consumerGroup = "chat-consumer-group"
	dbDSN         = "host=root-pg-chat-1 port=5432 dbname=chat user=chat-user password=chat-password sslmode=disable"
)

type serviceProvider struct {
	dbClient  db.Client
	txManager db.TxManager

	chatRepository repository.ChatRepository
	chatService    service.ChatService
	logRepository  repository.LogRepository

	authClient rpc.AuthClient
	apiClient  clientApi.ApiCat

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

func (s *serviceProvider) AuthClient() rpc.AuthClient {
	if s.authClient == nil {
		creds, err := credentials.NewClientTLSFromFile("certs/service.pem", "")
		if err != nil {
			log.Fatalf("failed to get credentials of authentication service: %v", err)
		}
		authConn, err := grpc.NewClient(fmt.Sprintf(":%d", servicePort),
			grpc.WithTransportCredentials(creds),
		)
		if err != nil {
			log.Fatalf("failed to connect to auth service: %v", err)
		}

		s.authClient = rpcAuth.NewAuthClient(access_v1.NewAccessV1Client(authConn))
	}
	return s.authClient
}

func (s *serviceProvider) ApiService() clientApi.ApiCat {
	if s.apiClient == nil {
		httpClient := &http.Client{
			Timeout: 3 * time.Second,
		}
		s.apiClient = apiHttp.NewApiClient(httpClient)
	}
	return s.apiClient
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
			s.ApiService(),
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

func (s *serviceProvider) KafkaConsumer() *consumer.Consumer {
	h := handler.NewHandler()
	c, err := consumer.NewConsumer(h, kafkaAddresses, topic, consumerGroup)
	if err != nil {
		log.Fatalf("failed to create kafka consumer: %v", err)
	}

	return c
}
