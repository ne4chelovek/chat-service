package tests

import (
	"context"
	"fmt"
	"github.com/ne4chelovek/chat_common/pkg/db"
	dbMocks "github.com/ne4chelovek/chat_common/pkg/db/mocks"
	"github.com/ne4chelovek/chat_common/pkg/db/transaction"
	clientApi "github.com/ne4chelovek/chat_service/internal/openApi"
	"github.com/ne4chelovek/chat_service/internal/repository"
	repoMocks "github.com/ne4chelovek/chat_service/internal/repository/mocks"
	"github.com/ne4chelovek/chat_service/internal/service/chat"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/brianvoe/gofakeit"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository
	type logRepositoryMockFunc func(mc *minimock.Controller) repository.LogRepository
	type apiInternalMock func(mc *minimock.Controller) clientApi.ApiCat
	type transactorMockFunc func(mc *minimock.Controller) db.Transactor

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		opts = pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

		repositoryErr = fmt.Errorf("failed to delete chat")

		req = id

		log = fmt.Sprintf("deleted chat with id: %v", id)
	)

	test := []struct {
		name string
		args args
		want *emptypb.Empty
		err  error

		chatServiceMock   chatRepositoryMockFunc
		logRepositoryMock logRepositoryMockFunc
		apiInternalMock   apiInternalMock
		transactorMock    transactorMockFunc
	}{
		{
			name: "chat repository success delete case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.DeleteChatMock.Expect(minimock.AnyContext, id).Return(nil, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				mock.LogMock.Expect(minimock.AnyContext, log).Return(nil)
				return mock
			},
			transactorMock: func(mc *minimock.Controller) db.Transactor {
				mock := dbMocks.NewTransactorMock(mc)
				txMock := dbMocks.NewTxMock(mc)
				mock.BeginTxMock.Expect(minimock.AnyContext, opts).Return(txMock, nil)
				txMock.CommitMock.Expect(minimock.AnyContext).Return(nil)
				return mock
			},
		},
		{
			name: "chat repository delete error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  repositoryErr,
			chatServiceMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.DeleteChatMock.Expect(minimock.AnyContext, id).Return(nil, repositoryErr)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				return mock
			},
			transactorMock: func(mc *minimock.Controller) db.Transactor {
				mock := dbMocks.NewTransactorMock(mc)
				txMock := dbMocks.NewTxMock(mc)
				mock.BeginTxMock.Expect(minimock.AnyContext, opts).Return(txMock, nil)
				txMock.RollbackMock.Expect(minimock.AnyContext).Return(nil)
				return mock
			},
		},
		{
			name: "chat repository log delete error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  repositoryErr,
			chatServiceMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repoMocks.NewChatRepositoryMock(mc)
				mock.DeleteChatMock.Expect(minimock.AnyContext, id).Return(nil, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repoMocks.NewLogRepositoryMock(mc)
				mock.LogMock.Expect(minimock.AnyContext, log).Return(repositoryErr)
				return mock
			},
			transactorMock: func(mc *minimock.Controller) db.Transactor {
				mock := dbMocks.NewTransactorMock(mc)
				txMock := dbMocks.NewTxMock(mc)
				mock.BeginTxMock.Expect(minimock.AnyContext, opts).Return(txMock, nil)
				txMock.RollbackMock.Expect(minimock.AnyContext).Return(nil)
				return mock
			},
		},
	}

	for _, tt := range test {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			chatServiceMock := tt.chatServiceMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			clientApiMock := tt.apiInternalMock(mc)
			transactorMock := transaction.NewTransactionManager(tt.transactorMock(mc))
			service := chat.NewService(chatServiceMock, logRepositoryMock, clientApiMock, transactorMock)

			res, err := service.DeleteChat(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, res)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
